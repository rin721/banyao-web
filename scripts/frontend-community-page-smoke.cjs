const { spawn, spawnSync } = require("child_process")
const fs = require("fs")
const path = require("path")

const repoRoot = path.resolve(__dirname, "..")
const backendPath = path.join(repoRoot, "backend")
const frontendPath = path.join(repoRoot, "frontend")
const workPath = path.join(repoRoot, "tmp", "ai", "frontend-community-page-smoke")
const screenshotsPath = path.join(workPath, "screenshots")
const backendBinaryPath = path.join(workPath, "console-server.exe")
const backendConfigPath = path.join(backendPath, "configs", "config.example.yaml")
const nuxtEntryPath = path.join(frontendPath, "node_modules", "nuxt", "bin", "nuxt.mjs")
const playwrightModulePath = path.join(repoRoot, "backend", "web", "app", "node_modules", "@playwright", "test")

const options = parseArgs(process.argv.slice(2))
const backendPort = Number(options["backend-port"] || 19997)
const frontendPort = Number(options["frontend-port"] || 3001)
const timeoutMs = Number(options["timeout-seconds"] || 90) * 1000
const backendBaseUrl = `http://127.0.0.1:${backendPort}`
const communityApiBaseUrl = `${backendBaseUrl}/api/v1/public/community`
const frontendBaseUrl = `http://127.0.0.1:${frontendPort}`

let backendProcess = null
let frontendProcess = null

main().catch(async (error) => {
  console.error(error.stack || error.message || String(error))
  await shutdown()
  cleanupBulkyArtifacts()
  process.exit(1)
})

async function main() {
  ensurePrerequisites()
  fs.mkdirSync(screenshotsPath, { recursive: true })

  console.log("Building backend community smoke server...")
  const build = spawnSync("go", ["build", "-mod=readonly", "-o", backendBinaryPath, "./cmd/console"], {
    cwd: backendPath,
    env: normalizedEnv(),
    stdio: "inherit",
    windowsHide: true
  })
  if (build.status !== 0) {
    throw new Error(`go build failed with exit code ${build.status}`)
  }

  backendProcess = startProcess(backendBinaryPath, ["server", `--config=${backendConfigPath}`], {
    cwd: backendPath,
    env: normalizedEnv({
      APP_SERVER_PORT: String(backendPort),
      APP_DB_DRIVER: "sqlite",
      APP_DB_SQLITE_PATH: path.join(workPath, "app.db"),
      APP_STORAGE_DRIVER: "local",
      APP_STORAGE_LOCAL_BASE_PATH: path.join(workPath, "uploads"),
      APP_LOG_FILE_PATH: path.join(workPath, "app.log"),
      APP_CORS_ALLOW_ORIGINS: frontendBaseUrl,
      APP_CORS_ALLOW_CREDENTIALS: "true",
      APP_AUTH_SIGNING_KEY: "frontend-page-smoke-signing-key-32",
      APP_AUTH_REFRESH_TOKEN_PEPPER: "frontend-page-smoke-refresh-pepper-32",
      APP_AUTH_MFA_SECRET_KEY: "frontend-page-smoke-mfa-secret-key-32",
      AUTH_SIGNING_KEY: "frontend-page-smoke-signing-key-32",
      AUTH_REFRESH_TOKEN_PEPPER: "frontend-page-smoke-refresh-pepper-32",
      AUTH_MFA_SECRET_KEY: "frontend-page-smoke-mfa-secret-key-32"
    }),
    logFile: path.join(workPath, "backend.log")
  })
  await waitForJson(`${communityApiBaseUrl}/status`, (json) =>
    json.code === 0 &&
    json.data &&
    json.data.mode === "go" &&
    Array.isArray(json.data.endpoints) &&
    json.data.endpoints.includes("/home")
  )

  frontendProcess = startProcess(process.execPath, [nuxtEntryPath, "dev", "--host", "127.0.0.1", "--port", String(frontendPort)], {
    cwd: frontendPath,
    env: normalizedEnv({
      BROWSER: "none",
      CI: "1",
      NUXT_PUBLIC_API_BASE_URL: communityApiBaseUrl,
      NUXT_PUBLIC_AUTH_API_BASE_URL: `${backendBaseUrl}/api/v1`,
      NUXT_PUBLIC_API_MOCK: "false"
    }),
    logFile: path.join(workPath, "frontend.log")
  })
  await waitForHtml(frontendBaseUrl)

  const { chromium } = require(playwrightModulePath)
  const browser = await chromium.launch({ headless: true })
  try {
    const results = []
    for (const viewport of [
      { name: "desktop", width: 1440, height: 900 },
      { name: "mobile", width: 390, height: 844 }
    ]) {
      const page = await browser.newPage({ viewport })
      const consoleErrors = []
      const failedRequests = []

      page.on("console", (message) => {
        if (message.type() === "error") {
          consoleErrors.push(message.text())
        }
      })
      page.on("requestfailed", (request) => {
        const url = request.url()
        if (!url.includes("/__nuxt")) {
          failedRequests.push(`${request.method()} ${url}: ${request.failure()?.errorText || "failed"}`)
        }
      })

      const home = await checkHomePage(page, viewport)
      const category = await checkCategoryPage(page, viewport)
      await page.close()

      if (consoleErrors.length > 0) {
        throw new Error(`Browser console errors on ${viewport.name}: ${consoleErrors.join(" | ")}`)
      }
      if (failedRequests.length > 0) {
        throw new Error(`Failed requests on ${viewport.name}: ${failedRequests.join(" | ")}`)
      }

      results.push({ viewport: viewport.name, home, category })
    }

    for (const result of results) {
      console.log(`[${result.viewport}] home videos=${result.home.videoCards}, dynamics=${result.home.dynamicCards}; category cards=${result.category.categoryCards}, maxCardWidth=${result.category.maxCategoryCardWidth}px`)
    }
    console.log(`Frontend community page smoke passed. Screenshots: ${screenshotsPath}`)
  } finally {
    await browser.close()
    await shutdown()
    cleanupBulkyArtifacts()
  }
}

async function checkHomePage(page, viewport) {
  await page.goto(frontendBaseUrl, { waitUntil: "networkidle" })
  await page.waitForSelector(".brand-band", { timeout: timeoutMs })
  await page.waitForSelector(".video-card", { timeout: timeoutMs })
  await page.screenshot({ path: path.join(screenshotsPath, `home-${viewport.name}.png`), fullPage: true })

  return await page.evaluate(() => {
    const brandBand = document.querySelector(".brand-band")
    const brandInner = document.querySelector(".brand-band__inner")
    const readSurface = (element) => {
      const style = window.getComputedStyle(element)
      return {
        backgroundColor: style.backgroundColor,
        borderTopWidth: style.borderTopWidth,
        boxShadow: style.boxShadow
      }
    }
    const home = {
      dynamicCards: document.querySelectorAll(".community-pulse__card").length,
      hero: readSurface(brandBand),
      heroInner: readSurface(brandInner),
      videoCards: document.querySelectorAll(".video-card").length
    }

    if (document.querySelector(".page-state__title")?.textContent?.includes("失败")) {
      throw new Error("Home page rendered an API failure state")
    }
    if (home.videoCards < 4 || home.dynamicCards < 1) {
      throw new Error(`Home page did not render backend feed data: ${JSON.stringify(home)}`)
    }
    if (
      home.hero.backgroundColor !== "rgba(0, 0, 0, 0)" ||
      home.hero.borderTopWidth !== "0px" ||
      home.hero.boxShadow !== "none" ||
      home.heroInner.backgroundColor !== "rgba(0, 0, 0, 0)" ||
      home.heroInner.borderTopWidth !== "0px" ||
      home.heroInner.boxShadow !== "none"
    ) {
      throw new Error(`Home brand band surface is not transparent: ${JSON.stringify(home)}`)
    }

    return home
  })
}

async function checkCategoryPage(page, viewport) {
  await page.goto(`${frontendBaseUrl}/category`, { waitUntil: "networkidle" })
  await page.waitForSelector(".category-card", { timeout: timeoutMs })
  await page.screenshot({ path: path.join(screenshotsPath, `category-${viewport.name}.png`), fullPage: true })

  return await page.evaluate(() => {
    const text = document.body.innerText
    const cards = Array.from(document.querySelectorAll(".category-card")).map((element) => {
      const box = element.getBoundingClientRect()
      return {
        height: Math.round(box.height),
        width: Math.round(box.width)
      }
    })
    const maxCategoryCardWidth = cards.reduce((max, item) => Math.max(max, item.width), 0)
    const category = {
      categoryCards: cards.length,
      hasBackendCategory: text.includes("创作") || text.includes("知识"),
      maxCategoryCardWidth
    }

    if (document.querySelector(".page-state__title")?.textContent?.includes("失败")) {
      throw new Error("Category page rendered an API failure state")
    }
    if (category.categoryCards < 5 || !category.hasBackendCategory) {
      throw new Error(`Category page did not render backend category data: ${JSON.stringify(category)}`)
    }
    if (window.innerWidth >= 640 && maxCategoryCardWidth > 270) {
      throw new Error(`Category cards are too wide for the compact map layout: ${JSON.stringify({ cards, maxCategoryCardWidth })}`)
    }

    return category
  })
}

function startProcess(command, args, options) {
  const output = fs.createWriteStream(options.logFile, { flags: "w" })
  const child = spawn(command, args, {
    cwd: options.cwd,
    env: options.env,
    stdio: ["ignore", "pipe", "pipe"],
    windowsHide: true
  })

  child.stdout.pipe(output)
  child.stderr.pipe(output)
  child.on("exit", (code, signal) => {
    if (code !== 0 && signal !== "SIGTERM") {
      console.error(`${path.basename(command)} exited with code ${code}`)
    }
  })

  return { child, output, logFile: options.logFile }
}

async function waitForJson(url, validate) {
  const response = await waitForResponse(url)
  const json = await response.json()
  if (!validate(json)) {
    throw new Error(`Unexpected JSON response from ${url}: ${JSON.stringify(json)}`)
  }
  return json
}

async function waitForHtml(url) {
  await waitForResponse(url, (response) => response.ok && String(response.headers.get("content-type") || "").includes("text/html"))
}

async function waitForResponse(url, validate = (response) => response.ok) {
  const startedAt = Date.now()
  let lastError = ""
  while (Date.now() - startedAt < timeoutMs) {
    for (const runtime of [backendProcess, frontendProcess]) {
      if (runtime && runtime.child.exitCode !== null) {
        const log = readTail(runtime.logFile)
        throw new Error(`Process for ${runtime.logFile} exited with ${runtime.child.exitCode}.\n${log}`)
      }
    }
    try {
      const response = await fetch(url)
      if (validate(response)) {
        return response
      }
      lastError = `HTTP ${response.status}`
    } catch (error) {
      lastError = error.message
    }
    await delay(500)
  }
  throw new Error(`Timed out waiting for ${url}: ${lastError}`)
}

async function shutdown() {
  for (const runtime of [frontendProcess, backendProcess]) {
    if (!runtime) {
      continue
    }
    if (runtime.child.exitCode === null) {
      runtime.child.kill()
      await Promise.race([
        onceExit(runtime.child),
        delay(5000)
      ])
      if (runtime.child.exitCode === null) {
        runtime.child.kill("SIGKILL")
      }
    }
    runtime.output.end()
  }
}

function cleanupBulkyArtifacts() {
  for (const target of [
    backendBinaryPath,
    path.join(workPath, "app.db"),
    path.join(workPath, "uploads"),
    path.join(workPath, "app.log")
  ]) {
    fs.rmSync(target, { force: true, recursive: true })
  }
}

function onceExit(child) {
  return new Promise((resolve) => {
    child.once("exit", resolve)
  })
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

function ensurePrerequisites() {
  for (const file of [backendConfigPath, nuxtEntryPath]) {
    if (!fs.existsSync(file)) {
      throw new Error(`Required file not found: ${file}`)
    }
  }
  if (!fs.existsSync(path.join(playwrightModulePath, "index.js"))) {
    throw new Error(`Playwright dependency not found: ${playwrightModulePath}`)
  }
}

function normalizedEnv(extra = {}) {
  const env = { ...process.env, ...extra }
  const pathValue = env.Path || env.PATH || ""
  delete env.PATH
  env.Path = pathValue
  return env
}

function readTail(filePath) {
  try {
    const content = fs.readFileSync(filePath, "utf8")
    return content.split(/\r?\n/).slice(-80).join("\n")
  } catch {
    return ""
  }
}

function parseArgs(args) {
  const out = {}
  for (let index = 0; index < args.length; index += 1) {
    const key = args[index]
    if (!key.startsWith("--")) {
      continue
    }
    out[key.slice(2)] = args[index + 1]
    index += 1
  }
  return out
}
