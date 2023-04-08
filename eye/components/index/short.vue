<i18n>
{
  "en": {
    "digital_artist": "Digital Artist",
    "blockquote_chatgpt": "Each artwork is not merely a combination of colors, shapes, and lines; it has its own soul and meaning.<br/><span class=\"is-size-7\">*ChatGPT about Artchitect</span>",
    "digital_creator": "Artchitect is a digital creator capable of producing paintings without any human involvement.",
    "autonomy": "Artchitect's autonomy relies on a quantum random number generator, which allows it to use universe light noise to create sequences of numbers. These numbers give birth to each painting's concept, plot, colors, style, and a distinct version. Origins of these paintings (quantum ideas) are stored at https://artchitect.space alongside each artwork.",
    "creations": "Artchitect's creations are made from the exceptional material of quantum randomness and complexity of the Universe. It can produce more than a thousand paintings every day, with one painting taking roughly one minute to create. Although each picture can be replicated indefinitely, their source - quantum ideas remains unique and impossible to replicate. These ideas appear only once, at the right time and in the right place.",
    "sd15": "All of the images are created using Stable Diffusion v1.5."
  },
  "ru": {
    "digital_artist": "Цифровой творец",
    "blockquote_chatgpt": "Произведение искусства - это не просто сочетание красок, форм и линий, это - целое, имеющее свою душу и значение.<br/><span class=\"is-size-7\">*ChatGPT про Artchitect</span>",
    "digital_creator": "Архитектор - цифровой творец, способный создавать картины без участия людей.",
    "autonomy": "Основой самостоятельности Архитектора является квантовый генератор случайных чисел. Архитектор использует световой шум Вселенной, чтобы создавать последовательности чисел, которые становятся идеей, сюжетом, цветами, стилем и конкретным вариантом каждой картины. Первоисточники картин (цифровые слепки) хранятся в галерее artchitect.space рядом с каждой картиной.",
    "creations": "Artchitect создаёт свои творения из особого материала - квантовой случайности и сложности Вселенной. Он создает одну картину в минуту и более тысячу картин каждый день. И хотя каждая картина может быть реплицирована бесконечное число раз, их источник - квантовые идеи - остаются уникальными и неповторимыми. Они возникли только один раз в нужное время и в нужном месте.",
    "sd15": "Для создания всех картин используется Stable Diffusion v1.5."
  }
}
</i18n>
<template>
  <section class="content">
    <h1><img src="/icon64.png" alt="artchitect"/> Artchitect</h1>
    <p class="is-size-5">
      🖌 {{$t('digital_artist')}}
    </p>
    <blockquote v-html="$t('blockquote_chatgpt')"/>
    <p>
      {{$t('digital_creator')}}
    </p>
    <p>
      {{$t('autonomy')}}
    </p>
    <p>
      {{$t('creations')}}
    </p>
    <p>
      {{$t('sd15')}}
    </p>
  </section>
</template>
<script>
export default {}
</script>
