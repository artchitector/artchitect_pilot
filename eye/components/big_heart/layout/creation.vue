<template>
  <div ref="resultBox">
    <div v-if="!message || !message.CardID" class="heart-heading">
      <h1 class="is-size-6 has-text-success">currently dreaming</h1>
      <div v-if="message" class="mb-3">
        <div>Seed</div>
        <span class="tag is-primary">{{ message.Seed }}</span>
      </div>
      <div v-if="message">
        <div>Tags</div>
        <div class="tags mb-3">
          <span class="tag" v-for="tag in message.Tags">{{ tag }}</span>
        </div>
      </div>
      <div>dreaming</div>
      <progress class="progress is-primary" :value="progress" max="100">-</progress>
    </div>
    <div v-else class="heart-result">
      <div ref="preImgElement" class="mb-3">
        enjoy the
        <NuxtLink class="has-text-info" :to="localePath(`/dream/${message.CardID}`)">#{{ message.CardID }}</NuxtLink>
        <progress class="progress is-warning" :value="enjoy" max="100">-</progress>
      </div>
      <a v-if="sizePrepared" :href="`/api/image/f/${message.CardID}`">
        <img :src="`/api/image/f/${message.CardID}`" :style="{'max-height': `${maxImgHeight}px`}"/>
      </a>
    </div>
  </div>
</template>

<script>
export default {
  name: "creation",
  data() {
    return {
      message: null,
      sizePrepared: false,
      maxImgHeight: 0,
      resizeTimeout: null,
    }
  },
  computed: {
    progress() {
      if (!this.message) {
        return 0;
      }
      if (!this.message.LastCardPaintTime || !this.message.CurrentCardPaintTime) {
        return 0;
      }
      const progress = this.message.CurrentCardPaintTime / this.message.LastCardPaintTime;
      return Math.floor(progress * 100);
    },
    enjoy() {
      if (!this.message) {
        return 0;
      }
      if (!this.message.EnjoyTime || !this.message.CurrentEnjoyTime) {
        return 0;
      }
      const progress = this.message.CurrentEnjoyTime / this.message.EnjoyTime;
      return Math.floor(progress * 100);
    }
  },
  mounted() {
    window.addEventListener('resize', this.onResize)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.onResize)
  },
  methods: {
    onMessage(channelName, message) {
      if (channelName !== "creation") {
        this.message = null
        return
      }
      this.message = message
      if (!this.sizePrepared) {
        this.fixImgSize()
      }
    },
    fixImgSize() {
      if (!this.$refs.preImgElement) {
        return
      }
      const article = this.$refs.resultBox.parentElement.parentElement
      const maxImgHeight = article.clientHeight - this.$refs.preImgElement.clientHeight
      this.maxImgHeight = maxImgHeight - 36
      console.log(`heart calculate height ${this.maxImgHeight}`)
      this.sizePrepared = true
    },
    onResize() {
      this.sizePrepared = false
      clearTimeout(this.resizeTimeout)
      this.resizeTimeout = setTimeout(() => {
        this.fixImgSize()
      }, 50)
    }
  }
}
</script>

<style scoped>
.heart-heading {
  letter-spacing: 1px;
  margin-bottom: 5px;
  text-transform: uppercase;
  font-size: 10px;
  padding: 10px;
  text-align: center;
}

.heart-result {
  letter-spacing: 1px;
  margin-bottom: 5px;
  font-size: 12px;
  text-transform: uppercase;
  padding: 10px;
  text-align: center;
}
</style>
