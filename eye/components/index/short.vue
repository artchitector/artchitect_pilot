<i18n>
{
  "en": {
    "heading_part": "is digital AI-artist who creates artworks from light",
    "middle_first": "💡 Light is a source of inspiration for an Artchitect. The idea of each card is formed from the beams of light that fall into the webcam. ",
    "middle_second": "Camera picks up the light, turns it into digital signals, signals into numbers, numbers into words, and words and numbers into pictures, and so on every 60 seconds.",
    "autonomy": "🤖 Artchitect is autonomous. It use the light to make ideas for cards without human help.",
    "ai": "🧠 Artchitect creates paintings with the help of modern artificial intelligence <a href=\"https://github.com/Stability-AI/stablediffusion\" target=\"_blank\">Stable Diffusion v1.5</a> (powered by <a href=\"https://github.com/invoke-ai/InvokeAI\" target=\"_blank\">Invoke.AI</a>). Source code available on <a href=\"https://github.com/artchitector/artchitect\" target=\"_blank\">github.com</a>.",
    "license": "🤝 All the paintings on this site belong to humanity. You can use any paintings from the artchitect.space website for any purpose, including commercial."
  },
  "ru": {
    "heading_part": " - это цифровой ИИ-художник, превращающий свет в картины.",
    "middle_first": "💡 Свет - источник вдохновения для Архитектора. Идея каждой картины сформирована из лучей света, которые падают в глазок веб-камеры.",
    "middle_second": "Камера улавливает кванты света, превращает это в цифровые сигналы, сигналы - в числа, числа - в слова, а слова и числа - в картины, и так каждые 60 секунд.",
    "autonomy": "🤖 Архитектор автономен. Он сам с помощью света придумывает идеи для картин без помощи человека.",
    "ai": "🧠 Архитектор рисует свои картины с помощью современного искусственного интеллекта <a href=\"https://github.com/Stability-AI/stablediffusion\" target=\"_blank\">Stable Diffusion v1.5</a> (powered by <a href=\"https://github.com/invoke-ai/InvokeAI\" target=\"_blank\">Invoke.AI</a>). Весь исходный код Архитектора доступен на <a href=\"https://github.com/artchitector/artchitect\" target=\"_blank\">github.com</a>.",
    "license": "🤝 Все картины на этом сайте принадлежат человечеству. Вы можете использовать любые картины с сайта artchitect.space в любых целях, в том числе коммерческих."
  }
}
</i18n>
<template>
  <section class="content">
    <h1><img src="/icon64.png" alt="artchitect"/> Artchitect</h1>
    <p class="is-size-5">
      🖌 <b>Artchitect</b> {{ $t('heading_part') }}
    </p>
    <p>
      {{ $t('middle_first') }}
      {{ $t('middle_second') }}
    </p>
    <p v-html="$t('autonomy')"/>
    <p v-html="$t('ai')"/>
    <p v-html="$t('license')"/>
    <p class="is-size-7 has-text-right">
      powered by <a href="https://github.com/invoke-ai/InvokeAI" target="_blank">Invoke.AI</a> (Stable Diffusion v1.5)
    </p>
  </section>
</template>
<script>
export default {}
</script>
