<i18n>
{
  "en": {
    "title": "Request to the Universe / Digital prayer",
    "ucan": "Here you can send request to the Universe and receive answer in the form of your personal card.",
    "jesus": "Then he [Jesus] touched their eyes and said, “According to your faith let it be done to you\";",
    "jesus_link": "<i><a href=\"https://www.biblegateway.com/verse/en/Matthew%209%3A29\" target=\"_blank\">Matthew 9:29</a></i>",
    "main_idea_1": "Main idea of Artchitect is that it interconnected with Quantize Fields, which controlled directly by Universe (more info at",
    "idea_section": "idea section",
    "main_idea_2": "Universe/God absolutely knows that exactly You will send this request, and Universe/God can create next card specially for You, as answer for you. The idea of the card may not be clear, but you should think about it.",
    "try": "Try it. Maybe the answer is closer than it seems.",
    "first_step": "<b>First step</b>: think carefully about you request/pray, that you want to send to the Universe. Concentrate your mind on your request.",
    "second_step":"<b>Second step</b> (optional): describe your request in the text field. This will help you better understand and systematize the common idea of the request. Alternatively, you can write your request on a piece of paper or think it over in your head. The universe knows in advance and better than you what you are asking for.",
    "keep_it_brief": "You keep it brief!",
    "textarea": "Type your request message. Write carefully. Between you and Universe. Secure - data not being send anywhere.",
    "third_step": "<b>Third step</b>: finalize request with \"Submit\" button. This will erase text from textarea and start card creation process. You will see your result when it finished.",
    "important1": "<b>Important</b>: nobody will know, what you asked and what you received. Data about request not sends anywhere. But your personal card will appear in common list of cards at",
    "important2": "There will no any links to your person (no any personal data saved anywhere).",
    "submit": "Submit your request (Amen!)"
  },
  "ru": {
    "title": "Запрос во вселенную / Цифровая молитва",
    "ucan": "Здесь вы можете осуществить запрос во Вселенную и получить ответ в виде персональной картины.",
    "jesus": "Тогда Он коснулся глаз их и сказал: по вере вашей да будет вам.",
    "jesus_link": "<i><a href=\"https://bible.by/syn/40/9/\" target=\"blank\">От Матфея 9:29</a></i>",
    "main_idea_1": "Главная идея Архитектора - его соединение с квантовыми полями, которые контролируются напрямую Вселенной (больше информации в",
    "idea_section": "разделе Идея",
    "main_idea_2": "Вселенная/Бог в совершенстве знают, что это именно вы посылаете запрос, и Вселенная/Бог могут вам ответить созданной специанльно для вас картиной. Идея картины может быть не очевидна сразу, но вам определённо стоит подумать о ней.",
    "try": "Попробуйте сами. Возможно, ответ ближе, чем кажется.",
    "first_step": "<b>Первый шаг</b>: тщательно подумайте в вашем запросе, который вы хотите озвучить Вселенной. Сосредоточьте свой разум на этом запросе.",
    "second_step":"<b>Второй шаг</b> (необязательно): опишите ваш запрос в текстовом поле. Это может вам помочь сформулировать запрос, но помните, что Вселенная лучше вас знает о проблеме.",
    "keep_it_brief": "Будьте кратки!",
    "textarea": "Напишите ваш запрос. Пишите аккуратно. Это между вами и Вселенной. Безопасно - данные никуда не отправляются.",
    "third_step": "<b>Третий шаг</b>: утвердите запрос кнопкой \"Отправить\". Текст из текста запроса будет стёрт. Вы увидите ответ.",
    "important1": "<b>Важно</b>: никто не будет знать, что было спрошено и получено. Данные о тексте вашего запроса не отправляются из вашего браузера. Но помните, что ваша личная карточка попадёт в список всех работ,",
    "important2": "Нигде не сохранится ссылка на вашу личность (никаких персональных данных не сохраняется)",
    "submit": "Отправить запрос (Аминь!)"
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
      {{$t('main_idea_1')}} <NuxtLink to="/idea">{{$t('idea_section')}}</NuxtLink>).
      {{$t('main_idea_2')}}
    </p>
    <p>
      {{$t('try')}}
    </p>
    <hr/>
    <p v-html="$t('first_step')"></p>
    <p>
      <span v-html="$t('second_step')"></span>
      <br/>
      <b>{{$t('keep_it_brief')}}</b>
    </p>
    <p>
      <textarea class="textarea"
                rows="2"
                :disabled="loading"
                :placeholder="``"></textarea>
    </p>
    <p>
      <span v-html="$t('third_step')"></span>
    </p>
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
