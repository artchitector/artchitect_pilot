<template>
  <div class="image-container">
    <NuxtLink :to="localePath(`/search/${hundred.Rank}/${hundred.Hundred}`)">
      <img :src="`/api/image_hundred/m/${hundred.Rank}/${hundred.Hundred}`"/>
    </NuxtLink>
    <p class="hundred-title">{{ hundred.Hundred }}-{{ hundred.Hundred + hundred.Rank - 1 }}</p>
  </div>
</template>

<script>
export default {
  name: "hundred",
  props: ['hundred']
}
</script>

<style lang="scss" scoped>
.image-container {
  position: relative;

  p.hundred-title {
    position: absolute;
    margin: auto;
    background-color: rgba(0, 0, 0, 0.6);
    padding: 5px;
    top: 30px;
    width: 100%;
    text-align: center;
    color: #d4d1c3;
  }
}
</style>
