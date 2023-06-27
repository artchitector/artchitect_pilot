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
    </div>
    <div v-else class="heart-result">
      enjoy the <a class="has-text-info" :href="`http://localhost/card/${message.CardID}`">#{{ message.CardID }}</a>
      <progress class="progress is-warning" :value="enjoy" max="100">-</progress>
      <a :href="`http://localhost/card/${message.CardID}`">
        <img :src="`/api/image/m/${message.CardID}`"/>
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
  methods: {
    onMessage(channelName, message) {
      this.message = message
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
