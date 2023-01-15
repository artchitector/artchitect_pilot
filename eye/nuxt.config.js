export default {
  // Global page headers: https://go.nuxtjs.dev/config-head
  head: {
    title: 'eye',
    htmlAttrs: {
      lang: 'en'
    },
    meta: [
      {charset: 'utf-8'},
      {name: 'viewport', content: 'width=device-width, initial-scale=1'},
      {hid: 'description', name: 'description', content: ''},
      {name: 'format-detection', content: 'telephone=no'}
    ],
    link: [
      {rel: 'icon', type: 'image/x-icon', href: '/icon64.png'}
    ]
  },

  // Global CSS: https://go.nuxtjs.dev/config-css
  css: [
    '~/assets/style.scss'
  ],

  // Plugins to run before rendering page: https://go.nuxtjs.dev/config-plugins
  plugins: [],

  // Auto import components: https://go.nuxtjs.dev/config-components
  components: true,

  // Modules for dev and build (recommended): https://go.nuxtjs.dev/config-modules
  buildModules: [
    '@nuxtjs/dotenv'
  ],

  // Modules: https://go.nuxtjs.dev/config-modules
  modules: [
    // https://go.nuxtjs.dev/buefy
    'nuxt-buefy',
    // https://go.nuxtjs.dev/axios
    '@nuxtjs/axios',
    // https://i18n.nuxtjs.org/setup
    '@nuxtjs/i18n'
  ],

  // Axios module configuration: https://go.nuxtjs.dev/config-axios
  axios: {
    // Workaround to avoid enforcing hard-coded localhost:3000: https://github.com/nuxt-community/axios-module/issues/308
    baseURL: process.env.SERVER_API_URL,
    browserBaseURL: process.env.CLIENT_API_URL
  },

  i18n: {
    /* module options */
    locales: ["en", "ru"],
    defaultLocale: 'ru',
    vueI18n: {
      fallbackLocale: 'en',
      messages: {
        en: {
          main: 'main',
          lottery: 'lottery',
          selection: 'selection',
          launched: ' launched 15th Jan 2023!',
          to_pray: 'Go to pray',
          loading: 'loading',
          last: 'last',
          cards: 'cards',
          page: 'page',
          lottery_description: 'Every day God chooses 10-100 cards out of ~1870 created in a day. He use lottery to chose.',
          selection_description: 'Every day God chooses 10-100 cards out of ~1870 created in a day. They go to selection.',
          card: 'Card',
          textarea_placeholder: 'Type your pray message. Write carefully. Between you and God. Secure, data burns and not being send anywhere.',
          pray_place: 'Pray place',
          wish: 'Wish one personal card as God\'s reply',
          burn: 'Pray! (burn text)',
          good_time_for_pray: 'Your answer loading. Good time for pray!',
          attempt: 'Attempt',
          usually_time: 'Usually that takes 2 minutes, but if artchitect is very loaded now, you need try once later.',
          something_wrong: 'Error',
          answer: 'Answer',
          to_cards: 'To cards'
        },
        ru: {
          main: 'главная',
          lottery: 'лотерея',
          selection: 'выбор',
          launched: ' открыт 15 января 2023 года!',
          to_pray: 'К молитве',
          loading: 'пожалуйста, подождите',
          last: 'последние',
          cards: 'карточек',
          page: 'страница',
          lottery_description: 'Каждый день Бог выбирает 10-100 карточек из ~1870, созданных за день. Для этого Он использует лотерею.',
          selection_description: 'Каждый день Бог выбирает 10-100 карточек из ~1870, созданных за день. Они попадают в "выбор".',
          card: 'Карточка',
          textarea_placeholder: 'Напишите текст молитвы или просьбы. Пишите вдумчиво. Это между вами с Богом. Безопасно, данные никуда не отправляются и сгорают',
          pray_place: 'Место для молитвы',
          wish: 'Пожелать одну персональную карточку как ответ/послание от Бога',
          burn: 'Помолиться! (сжечь записку)',
          good_time_for_pray: 'Ваш ответ рисуется. Сейчас хорошее время для молитвы!',
          attept: 'Попытка',
          usually_time: 'Обычно ожидание занимает 2 минуты, но если artchitect очень загружен, то попробуйте в другой раз. Бог вас всё равно слышит',
          something_wrong: 'Что-то случилось',
          answer: 'Ответ',
          to_cards: 'К картинам'
        }
      }
    }
  }
}
