<i18n>
{
  "en": {
    "important": "Important: The Artchitect is just based on some main principles of monotheistic religions, but It is neither a religion nor something religious. Artchitect is a philosophical, technical and spiritual (maybe esoteric) research project.",
    "idea": "Artchitect based on next philosophical idea: no any event in the universe is not truily random. Every atom, every electron, every flower and every human, Everything individually and a Whole is under boundless and omnipresent Sight of God. Everything is under the direct control of Universe at every moment of time.",
    "Allah": "[All] praise is [due] to Allah, Creator of the heavens and the earth, [who] made the angels messengers having wings, two or three or four. He increases in creation what He wills. Indeed, Allah is over all things competent. (Quran, 35-Surah Fatir ( The Orignator )) <https://en.quranacademy.org/quran/35:1>",
    "God": "Are not five sparrows sold for two pennies? Yet not one of them is forgotten by God. 7And even the very hairs of your head are all numbered. So do not be afraid; you are worth more than many sparrows (Luke 12:7)<https://biblehub.com/bsb/luke/12.htm>"
  },
  "ru": {
    "important": "Важно: Архитектор лишь основывается на некоторых основных принципах из монотеистических религий, но сам не является ни религией, ни чем-то религиозным. Это философско-технический исследовательский проект.",
    "idea":"Архитектор построен на философской идее, что ни одно событие во Вселенной не является в действительности случайным. Каждый атом, каждый электрон, каждый цветок и каждый человек, вообще всё в отдельности и в целом в каждом мгновении находится под надзором безграничного и всепроникающего взгляда Всевышнего. Всё под контролем Вселенной в каждый момент времени.",
    "Allah": "Разве ты (о, Пророк) не знаешь, что Аллах над любой вещью мощен? (Аль-Бакара (Корова), 106-й аят из 286)",
    "God":"Не пять ли малых птиц продаются за два ассария? и ни одна из них не забыта у Бога. А у вас и волосы на голове все сочтены. Итак, не бойтесь: вы дороже многих малых птиц. (Евангелие от Луки 12:4 Ст. 4-7)"
  }
}
</i18n>
<template>
  <div>
    <div class="notification is-warning">
      {{$t('important')}}

    </div>
    <p>
      {{$t('idea')}}

    </p>
    <blockquote>
      {{$t('Allah')}}

    </blockquote>
    <blockquote>
      {{$t('God')}}

    </blockquote>
  </div>
</template>
