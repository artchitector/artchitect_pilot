<i18n>
{
  "en": {
    "title": "Artchitect - entropy"
  },
  "ru": {
    "title": "Artchitect - энтропия"
  }
}
</i18n>

<template>
  <section>
    <h1 class="is-size-3">Датчик энтропии</h1>
    <hr/>
    <p>
      Для рисования картин Архитектору приходится придумывать "что рисовать?" - придумывать набор ключевых слов и
      уникальный seed-номер. Для возможности выбора (принятия решений)
      Архитектор подключён к датчику энтропии (собранного из веб-камеры), который позволяет из светового шума получить
      float64-число. <b>float64-число</b> - фундамент каждого выбора в работе Artchitect, а в день их происходит десятки
      тысяч.
    </p>
    <hr/>
    <p>
      <b>Шаг 1. Исходный кадр (Source)</b>
      <br/>
      Энтропия (в форме светового шума) получается из разницы двух соседних кадров. Текущий кадр виден ниже, его
      содержание не имеет
      значения. Красным отмечена область, с которой будет сниматься фоновый шум.
    </p>
    <div class="has-text-centered">
      <img v-if="images.source !== null" :src="`data:image/jpeg;base64, ${images.source}`"
           alt="loading source stream"/>
    </div>
    <p>
      <b>Шаг 2. Шум (Noise)</b>
      <br/>
      Для экстракции светового шума используется вычитание кадров (из цветов текущего кадра вычитаются цвета
      предыдущего кадра). Разница между кадрами и есть световой шум, который дополнительно усиливается. В визуальном
      представлении шум выглядит следующим образом:
    </p>
    <div class="has-text-centered">
      <img v-if="images.noise !== null" :src="`data:image/jpeg;base64, ${images.noise}`"
           style="height: 256px; width: 256px;"
           alt="loading noise stream"/>
    </div>
    <p>
      <b>Шаг 3. Энтропия - сжатие и нормализация шума (Entropy)</b>
      <br/>
      Шум сжимается до области 8x8 пикселей и нормализуется (распределяется по шкале от 0 до
      255). Так и выглядит <b>базовая энтропия</b>, как её видит Архитектор.
    </p>
    <div class="has-text-centered">

      <img v-if="images.shrink !== null" :src="`data:image/png;base64, ${images.shrink}`"
           style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading shrink stream"/>
    </div>
    <p><b>Шаг 4. Инвертированная энтропия (Choice)</b>
      <br/>
      Чтобы Архитектор рисовал разнообразные картины, он при каждом выборе (ключевого слова или seed-номера) должен
      обеспечивать максимальное разнообразие (Архитектор должен выбирать очень разные слова из словаря каждый раз,
      нечасто повторяясь).
      <br/>
      Чтобы
      обеспечить такую селективность, энтропия инвертируется бинарно (каждый байт в рисунке энтропии зеркально
      отражается). Инвертированная энтропия становится очень случайной и позволяет равномерно справедливо выбирать слова
      из словаря.
    </p>
    <div class="has-text-centered">
      <img v-if="images.bytes !== null" :src="`data:image/png;base64, ${images.bytes}`"
           style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading bytes stream"/>
      <br/>
    </div>
    <p>
      <b>Шаг 5. Превращение изображения в число</b>
      <br/>
      Каждый пиксель изображения или включен, или выключен (значение его цвета ближе к красному или к чёрному). Каждый
      пиксель - это один бит из 64-битного целого числа (uint64). Картинка после бинарного преобразования становится
      uint64 числом на шкале от 0 до 18446744073709551615. Положение этого числа на шкале и есть точка выбора.
      Далее шкала превращается в float64-число, представляющее шкалу выбора - это дробное число от 0.0 до 1.0, где 0 -
      самый первый элемент в списке, 1.0 - последний элемент, 0.5 - примерно середина списка.
      <br/>
      Архитектору хватает нескольких float64-чисел, чтобы принять все решения относительно идеи новой картины.
    </p>
    <div class="has-text-centered">
      <span class="is-size-7" v-html="entropy.bytes ? entropy.bytes.match(/.{1,8}/g).join('<br/>') : '-'"></span><br/>
      <span>Выбранное значение по шкале от 0.0 до 1.0: <b>{{ entropy.float }}</b></span><br/>
    </div>

  </section>
</template>

<script>

import WsConnection from "~/utils/ws_connection";

export default {
  name: "entropy",
  data() {
    return {
      player: null,
      logPrefix: "🎆",
      status: {
        error: null,
        reconnecting: null,
      },
      maintenance: false,
      connection: null,
      images: {
        source: null,
        noise: null,
        shrink: null,
        bytes: null
      },
      entropy: {
        bytes: "",
        float: 0.0
      }
    }
  },
  head() {
    return {
      title: this.$t('title')
    }
  },
  mounted() {
    if (process.env.SOUL_MAINTENANCE === 'true') {
      this.maintenance = true
      return
    }
    this.connection = new WsConnection(process.env.WS_URL, this.logPrefix, ['entropy'], 100)
    this.connection.onmessage((channel, message) => {
      this.status.error = null
      this.status.reconnecting = null
      this.onMessage(channel, message)
    })
    this.connection.onerror((err) => {
      this.status.error = err
    })
    this.connection.onreconnecting((attempt, maxAttempts) => {
      console.log(`${this.logPrefix}: RECONNECTING ${attempt}/${maxAttempts}`)
      this.status.reconnecting = {attempt, maxAttempts}
    })
    this.connection.onopen(() => {
      this.status.reconnecting = null
      this.status.error = null
    })
    this.connection.onopen(() => {
      this.status.reconnecting = null
      this.status.error = null
      console.log(`${this.logPrefix}: connection established`)
    })
    this.connection.connect()
  },
  beforeDestroy() {
    if (!this.maintenance) {
      this.connection.close()
      this.connection = null
    }
  },
  methods: {
    onMessage(chan, msg) {
      if (!msg.Image) {
        return
      }

      if (msg.EntropyAnswer > 0) {
        this.entropy.bytes = msg.EntropyAnswerByte
        this.entropy.float = msg.EntropyAnswer
      }

      switch (msg.Phase) {
        case 'source':
          this.images.source = msg.Image
          break
        case 'noise':
          this.images.noise = msg.Image
          break
        case 'shrink':
          this.images.shrink = msg.Image
          break
        case 'bytes':
          this.images.bytes = msg.Image
          break
        default:
          console.warn(`new phase ${msg.Phase}`)
      }
    }
  }
}
</script>


<style scoped lang="scss">

</style>
