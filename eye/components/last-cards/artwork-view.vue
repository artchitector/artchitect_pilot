<template>
  <div class="image-container">
    <NuxtLink :to="`/card/${artwork.ID}`" target="_blank">
      <img :src="'/api/image/m/' + artwork.ID"/>
      <div class="version">
        <span class="tag is-primary is-light">{{ artwork.Version }}</span>
      </div>
      <p class="is-size-7 tags">{{ artwork.Spell.Tags }}</p>
      <p class="is-size-7 info">id={{ artwork.ID }}, <b>seed: {{ artwork.Spell.Seed }}</b></p>
    </NuxtLink>
  </div>
</template>
<script>
export default {
  props: ['artwork']
}
</script>
<style lang="scss">
.image-container {
  position: relative;
  min-height: 431px;
  .version {
    position: absolute;
    right: 5px;
    bottom: 10px;
    opacity: 70%;
  }
}

.image-container p {
  visibility: visible;
  position: absolute;
  margin: auto;
  word-break: break-all;
  background-color: rgba(0, 0, 0, 0.6);
  padding: 5px;
}

.image-container p.tags {
  bottom: 30px;
}

.image-container p.info {
  top: 30px;
  width: 100%;
  text-align: center;
}

</style>
