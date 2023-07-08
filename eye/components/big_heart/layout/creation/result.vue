<template>
  <div class="heart-result">
    <template v-if="message">
      <div ref="preImgElement" class="mb-3">
        enjoy the
        <NuxtLink class="has-text-info" :to="localePath(`/dream/${message.CardID}`)">#{{ message.CardID }}</NuxtLink>
        <progress class="progress is-warning" :value="enjoy" max="100">-</progress>
      </div>
      <div v-if="message && sizePrepared" class="image-wrapper">
        <a :href="`/api/image/f/${message.CardID}`">
          <img :src="`/api/image/f/${message.CardID}`" :style="{'max-height': `${maxImgHeight}px`}"/>
        </a>
        <div class="control-like">
          <font-awesome-icon v-if="liked && liked.error"
                             icon="fa-solid fa-triangle-exclamation"
                             :title="liked.error.message"/>
          <a v-else href="#" @click.prevent="like()">
            <font-awesome-icon v-if="!liked || !liked.liked" icon="fa-solid fa-heart" class="has-color-base"/>
            <font-awesome-icon v-else icon="fa-solid fa-heart" class="has-text-danger"/>
          </a>
        </div>
      </div>
    </template>
  </div>
</template>

<script>
export default {
  name: "result",
  data() {
    return {
      message: null,
      sizePrepared: false,
      resizeTimeout: null,
      maxImgHeight: 0,
      liked: {
        liked: false,
        error: null,
      }
    }
  },
  computed: {
    enjoy() {
      if (!this.message) {
        return 0;
      }
      if (!this.message.EnjoyTime || !this.message.CurrentEnjoyTime) {
        return 0;
      }
      const progress = this.message.CurrentEnjoyTime / this.message.EnjoyTime;
      return Math.floor(progress * 100);
    },
  },
  mounted() {
    setTimeout(() => {
      this.initLiked()
    }, 1000)
  },
  methods: {
    onMessage(message) {
      this.message = message
      if (!this.sizePrepared) {
        this.fixImgSize()
      }
    },
    fixImgSize() {
      if (!this.$refs.preImgElement) {
        return
      }
      const article = this.$refs.preImgElement.parentElement.parentElement.parentElement.parentElement
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
        this.lastDreamBoxHeightValue = null
      }, 50)
    },
    async initLiked() {
      console.log('initLiked', this.message)
      if (!this.message) {
        return
      }
      const cardID = this.message.CardID
      try {
        let like = await this.$axios.$get(`/liked/${cardID}`)
        this.liked.liked = like.Liked
      } catch (e) {
        console.error(e)
        this.liked.error = e
      }
    },
    async like() {
      try {
        let like = await this.$axios.$post("/like", {
          card_id: this.message.CardID,
        })
        this.$emit('liked', like)
        this.liked = {
          id: like.ID,
          liked: like.Liked,
        };
      } catch (e) {
        console.error(e)
        this.liked = {
          error: e
        };
      }

    },
  }
}
</script>

<style lang="scss" scoped>
.heart-result {
  letter-spacing: 1px;
  margin-bottom: 5px;
  font-size: 12px;
  text-transform: uppercase;
  padding: 10px;
  text-align: center;
}

.image-wrapper {
  position: relative;

  .control-like {
    position: absolute;
    left: 50%;
    bottom: 20%;
    z-index: 3;
    margin-left: -20px;
    font-size: 48px;
    opacity: 0.7;
    filter: drop-shadow(0px 0px 8px rgba(255, 0, 0, 0.6));
  }
}
</style>
