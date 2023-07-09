<template>
  <div class="result-wrapper">
    <p class="is-size-5 has-text-success has-text-centered">created dream is
      <NuxtLink :to="localePath(`/dream/${dream_id}`)" class="has-text-info">#{{ dream_id }}</NuxtLink>
    </p>
    <progress class="progress is-warning" :value="progress" max="100">-</progress>
    <NuxtLink :to="localePath(`/dream/${dream_id}`)" class="has-text-info">
      <img :src="`/api/image/f/${dream_id}`"/>
    </NuxtLink>
  </div>
</template>

<script>
export default {
  name: "result",
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

  img {
    max-height: 90%;
  }
}
</style>
