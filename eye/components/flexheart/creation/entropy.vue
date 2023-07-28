<template>
  <div class="has-text-centered">
    <span class="is-size-7">ENTROPY</span><br/>
    <img v-if="entropy.ImagesEncoded.entropy !== null" :src="`data:image/png;base64, ${entropy.ImagesEncoded.entropy}`"
         style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading entropy stream"/>
    <img v-if="entropy.ImagesEncoded.choice !== null" :src="`data:image/png;base64, ${entropy.ImagesEncoded.choice}`"
         style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading choice stream"/>
  </div>
</template>

<script>
export default {
  props: ['entropy']
}
</script>

<style scoped lang="scss">

</style>
