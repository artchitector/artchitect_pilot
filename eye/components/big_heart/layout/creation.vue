<template>
  <div>
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
      <last_dream v-if="message" :message="message" :heartState="heartState" :style="{'height': lastDreamBoxHeight}"/>
    </div>
    <result v-else ref="resultBox"/>
  </div>
</template>

<script>
import Last_dream from "@/components/big_heart/layout/creation/last_dream.vue";
import Result from "@/components/big_heart/layout/creation/result.vue";

export default {
  name: "creation",
  components: {Result, Last_dream},
  data() {
    return {
      message: null,
      heartState: null,
      sizePrepared: false,
      maxImgHeight: 0,
      lastDreamBoxHeightValue: null,
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

    lastDreamBoxHeight() {
      if (this.lastDreamBoxHeightValue !== null) {
        return this.lastDreamBoxHeightValue;
      } else {
        const screenHeight = window.innerHeight
        console.log(screenHeight)
        if (screenHeight > 900) {
          this.lastDreamBoxHeightValue = `${Math.floor(screenHeight - 450)}px`
          return this.lastDreamBoxHeightValue
        }

        this.lastDreamBoxHeightValue = '40%'
        return this.lastDreamBoxHeightValue
      }
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
      if (channelName === "heart") {
        this.updateHeartState(message)
        return
      }
      if (channelName !== "creation") {
        this.message = null
        return
      }
      this.message = message
      if (this.$refs.resultBox) {
        this.$refs.resultBox.onMessage(this.message)
      }
    },
    updateHeartState(msg) {
      this.heartState = msg
    },
    onResize() {
      this.lastDreamBoxHeightValue = null
      if (this.$refs.resultBox) {
        this.$refs.resultBox.onResize()
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.heart-heading {
  letter-spacing: 1px;
  margin-bottom: 5px;
  text-transform: uppercase;
  font-size: 10px;
  padding: 10px;
  text-align: center;

  .tags .tag {
    font-size: 9px;
  }
}
</style>
