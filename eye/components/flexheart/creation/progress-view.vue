<template>
  <div class="progress-view">
    <div v-if="message.CardID">
      already created #{{message.CardID}}
    </div>
    <div v-else class="heart-heading">
      <entropy v-if="entropy" :entropy="entropy"/>
      <h1 class="is-size-5 has-text-success has-text-centered">currently creating</h1>
      <div>
        <div class="tags mb-3">
          <span class="tag is-primary">seed={{message.Seed}}</span>
          <span class="tag" v-for="tag in message.Tags">{{ tag }}</span>
        </div>
      </div>
      <div class="is-size-7 has-text-centered">creating</div>
      <progress class="progress is-primary" :value="progress" max="100">-</progress>
    </div>
  </div>
</template>

<script>
import Entropy from "@/components/flexheart/creation/entropy.vue";
export default {
  name: "progress-view",
  components: {Entropy},
  props: ["message", "entropy"],
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
  }
}
</script>

<style lang="scss" scoped>
  .progress-view {
    max-width: 800px;
    .tags .tag {
      font-size: 10px;
    }
  }
</style>
