<i18n>
{
  "en": {
    "title": "Artchitect - lottery",
    "lottery_description": "Everyday Artchitect selects 10-100 chosen cards, and saves them into ",
    "selection": "selection"
  },
  "ru": {
    "title": "Artchitect - лотерея",
    "lottery_description": "Ежедневно Архитектор выбирает 10-100 карточек из общего числа, созданных за день (~1400 карточек), и сохраняет их в ",
    "selection": "отборное"
  }
}
</i18n>
<template>
  <div>
    <section>
      {{ $t('lottery_description') }} <NuxtLink :to="localePath('selection')">{{$t('selection')}}</NuxtLink>.
    </section>
    <section>
      <lottery-list/>
    </section>
  </div>
</template>
<script>
export default {
  head() {
    return {
      title: this.$t('title')
    }
  }
}
</script>
/
