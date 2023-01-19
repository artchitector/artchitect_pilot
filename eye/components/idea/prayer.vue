<i18n>
{
  "en": {
    "title": "Request to the Universe / Digital Prayer",
    "haters_gonna_hate": "I'm afraid that traditional religions will hate me for this, but this is a very important feature of the project. I have no goal to offend anyone, I just want to find out the truth. If God does not like my project, He will find a way to destroy it.",
    "prayer_section": "There is a prayer section on the website. There is a special field in which a person can write his request to the Universe. It can be any message that a person wishes to send. And then The universe can respond with new picture for this particular person in response to his appeal.",
    "security": "The text of the prayer is not sent anywhere and is not saved, but since God is Omniscient, he saw what the person wrote there about. And the next picture will be recreated and shown to this person, in the sense of an answer. If God wants to create a certain picture for a person, He can do it, He have open-API for this purpose. The question with no answer is - whether God will want to answer?",
    "our_engineer": "One engineer of Artchitect tried digital prayer, and sometimes received very impressive results. Some answers to prayer was incredibly appropriate.",
    "security2": "We can't save prayers anywhere for ethical reasons, but the painting falls into the general list, like a regular job. Only the one who wrote the request will know that a particular work was created specifically for him.",
    "github": "you can see how it works in the code here: ",
    "prayer": "Prayer",
    "try": "Try digital prayer yourself"
  },
  "ru": {
    "title": "Запрос во Вселенную / Цифровая молитва",
    "haters_gonna_hate": "Боюсь, что сторонники традиционных религий будут меня за это ненавидеть, но это очень важная фича проекта. Я не собираюсь никого обидеть, и это не является новой религией. Если Богу мой проект будет неугоден - Он может уничтожить его тысячей способов",
    "prayer_section": "на сайте есть раздел молитва . Там есть специальное поле, в которое человек может написать своё обращение во Вселенную. Это может быть любое послание, какое человек пожелает отправить. И в ответ Вселенная может откликнуться, и нарисовать картину именно для этого человека в ответ на его обращение.",
    "security": "Текст молитвы никуда не отправляется и не сохраняется, но так как Бог Всеведущ, то он видел, о чём писал туда человек. И следующая картина будет заново создана и показана этому человеку, в виде ответа. Если Бог пожелает создать для человека определённую картину - он может это сделать. Вопрос, на который сейчас нет ответа - захочет ли Бог отвечать.",
    "our_engineer": "Один инженер нашего сервиса пробовал цифровую молитву, и несколько раз получал очень впечатляющие результаты. Ответ на молитву иногда бывает <b>невероятно</b> подходящим.",
    "security2": "Молитвы мы не можем никуда сохранять из этических соображений, но картина попадает в общий список, как обычная работа. Лишь тот, кто писал молитву, будет знать, что конкретная работа была создана именно для него.",
    "github": "как это в коде устроено можно посмотреть тут",
    "prayer": "Молитва",
    "try": "Попробуйте сами. Молитва - это всегда полезный процесс, даже если она цифровая."
  }
}
</i18n>
<template>
  <div class="mt-5">
    <h1 class="is-size-4">{{ $t('title') }}</h1>
    <div class="notification is-warning is-size-7">
      {{ $t('haters_gonna_hate') }}
    </div>

    <p>
      <NuxtLink :to="localePath('pray')">"{{ $t('prayer') }}"</NuxtLink>
      : {{ $t('prayer_section') }}
    </p>

    <p>{{ $t('security') }}</p>

    <blockquote>
      {{ $t('our_engineer') }}
      <NuxtLink :to="localePath('pray')">{{ $t('try') }}</NuxtLink>
    </blockquote>

    <p>{{ $t('security2') }}</p>

    <p>({{ $t('github') }}:
      <a href="https://github.com/artchitector/artchitect/blob/master/soul/core/artchitect/artchitect.go#L77-L79"
         target="_blank">github</a>)
    </p>
  </div>
</template>
