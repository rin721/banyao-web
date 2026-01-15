<script>
  // No script logic needed for now
</script>

<section class="hero-container">
  <div class="content">
    <h1 class="title">半妖<span class="sub">(PicHub)</span></h1>
    <p class="slogan">半妖的旧址，这是一个二次元图库。</p>

    <button class="cta-button"> 待开发中 qwq... </button>
  </div>
</section>

<style>
  .hero-container {
    position: relative;
    z-index: 10;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    width: 100%;
    text-align: center;
  }

  .content {
    animation: fadeIn 1.5s ease-out;
  }

  .title {
    font-size: 5rem;
    font-weight: 800;
    letter-spacing: 0.1em;
    margin-bottom: 0.5rem;
    color: #fff;
    text-shadow: 0 0 20px rgba(121, 190, 255, 0.932);
  }

  .title .sub {
    font-size: 2rem;
    font-weight: 400;
    opacity: 0.8;
    margin-left: 10px;
  }

  .slogan {
    font-size: 1.2rem;
    color: var(--color-text-sub);
    margin-bottom: 3rem;
    letter-spacing: 0.2em;
    font-weight: 300;
    filter: blur(0.3px); /* Slightly dreamy blur */
  }

  .cta-button {
    background-color: var(--color-primary);
    color: white;
    font-size: 1.1rem;
    padding: 1rem 3rem;
    border-radius: 9999px; /* Pill shape */
    box-shadow: 0 0 15px rgba(0, 210, 211, 0.4);
    transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
    position: relative;
    overflow: hidden;
  }

  .cta-button:hover {
    background-color: var(--color-primary-hover);
    transform: scale(1.05);
    box-shadow: 0 0 25px rgba(0, 210, 211, 0.7);
  }

  .cta-button:active {
    transform: scale(0.95);
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  /* Mobile responsiveness */
  @media (max-width: 768px) {
    .title {
      font-size: 3rem;
    }
    .title .sub {
      font-size: 1.5rem;
      display: block;
      margin: 5px 0 0 0;
    }
  }
</style>
