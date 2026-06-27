import {
  getMockCommunitySession,
  mockCommunityAuthCookieName
} from "../../../../shared/mocks/auth"

export default defineEventHandler((event) => {
  const sessionId = getCookie(event, mockCommunityAuthCookieName)
  const session = sessionId ? getMockCommunitySession(sessionId) : null

  return session || null
})
