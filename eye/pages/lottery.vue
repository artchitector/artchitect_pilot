<i18n>
{
  "en": {
    "title": "Artchitect - lottery",
    "lottery_description": "Every day God chooses 10-100 cards out of ~1400 created in a day. He use lottery to chose."
  },
  "ru": {
    "title": "Artchitect - лотерея",
    "lottery_description": "Ежедневно Бог выбирает 10-100 отобранных карточек из общего числа ~1400 карточек, созданных за день. Он использует лотерею для этого."
  }
}
</i18n>
<template>
  <div>
    <section>
      {{ $t('lottery_description') }}
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
