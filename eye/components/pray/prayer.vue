<i18n>
{
  "en": {
    "title": "Request",
    "ucan": "Here you can receive personal card from Artchitect in form of picture.",
    "jesus": "Then he [Jesus] touched their eyes and said, “According to your faith let it be done to you\";",
    "jesus_link": "<i><a href=\"https://www.biblegateway.com/verse/en/Matthew%209%3A29\" target=\"_blank\">Matthew 9:29</a></i>",
    "try": "Try it yourself!",
    "important1": "<b>Important</b>: your personal answer will be saved as any another Artchitect's card.",
    "important2": "Personal data not saved",
    "submit": "Get personal card"

  },
  "ru": {
    "title": "Запрос",
    "ucan": "Здесь вы можете осуществить запрос к Архитектору и получить ответ в виде персональной картины. Это можно использовать как молитву, гадание, астрологический прогноз, метафорические карты и т.д.",
    "jesus": "Тогда Он коснулся глаз их и сказал: по вере вашей да будет вам.",
    "jesus_link": "<i><a href=\"https://bible.by/syn/40/9/\" target=\"blank\">От Матфея 9:29</a></i>",
    "try": "Попробуйте сами. Возможно, ответ ближе, чем кажется.",
    "important1": "<b>Важно</b>: помните, что ваша персональная карточка-ответ попадёт в список всех работ,",
    "important2": "Нигде не сохранится ссылка на вашу личность (никаких персональных данных не сохраняется)",
    "submit": "Получить карточку"
  }
}
</i18n>
<template>
  <div class="content">
    <h1 class="is-size-4">{{$t('title')}}</h1>
    <p>
      {{$t('ucan')}}
    </p>
    <blockquote>
      {{$t('jesus')}}
      <br/>
      <div v-html="$t('jesus_link')"/>
    </blockquote>
    <p>
      {{$t('try')}}
    </p>
    <hr/>
    <div class="notification">
      <span v-html="$t('important1')"></span>
      <NuxtLink to="/">artchitect.space</NuxtLink>.
      <span v-html="$t('important2')"></span>
    </div>
    <p class="has-text-centered">
      <button class="button is-primary" :disabled="loading" @click.prevent="submit()">{{$t('submit')}}
      </button>
    </p>
    <p v-if="error" class="notification is-danger">{{ error }}</p>
    <p v-if="loading" class="has-text-centered">
      <loader/>
    </p>
  </div>
</template>

<script>
export default {
  name: "prayer",
  data() {
    const randomPassword = (Math.random() + 1).toString(36).substring(7);
    return {
      password: randomPassword,
      loading: false,
      error: null,
    }
  },
  methods: {
    async submit() {
      try {
        this.loading = true
        localStorage.setItem("last_pray_password", this.password);  // only emitter of pray can see answer (temporary password needed)
        const result = await this.$axios.$post("/pray", {
          password: this.password
        })
        let url = this.localePath(`/prayer/${result}`)
        await this.$router.push(url)
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>

</style>
