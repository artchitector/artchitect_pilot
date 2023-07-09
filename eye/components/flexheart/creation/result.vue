<template>
  <div class="result-wrapper">
    <p class="is-size-5 has-text-success has-text-centered">created dream is
      <NuxtLink :to="localePath(`/dream/${dream_id}`)" class="has-text-info">#{{ dream_id }}</NuxtLink>
    </p>
    <progress class="progress is-warning" :value="progress" max="100">-</progress>
    <div class="image-container">
      <NuxtLink :to="localePath(`/dream/${dream_id}`)" class="has-text-info">
        <img :src="`/api/image/f/${dream_id}`"/>
      </NuxtLink>
      <div class="control-like">
        <liker :dream_id="dream_id"/>
      </div>
    </div>
  </div>
</template>

<script>
import Liker from "@/components/utils/liker.vue";

export default {
  name: "result",
  components: {Liker},
  props: ["dream_id", "totalEnjoyTime", "currentEnjoyTime"],
  computed: {
    progress() {
      if (!this.totalEnjoyTime || !this.currentEnjoyTime) {
        return 0;
      }
      const progress = this.currentEnjoyTime / this.totalEnjoyTime;
      return Math.floor(progress * 100);
    },
  }

}
</script>

<style lang="scss" scoped>
.result-wrapper {
  text-align: center;
  height: 100%;

  .image-container {
    height: 90%;
    position: relative;
    img {
      max-height: 100%;
    }
    .control-like {
      position: absolute;
      left: 50%;
      bottom: 10%;
      z-index: 3;
      transform: translate(-50%, -10%);
    }
  }
}
</style>
