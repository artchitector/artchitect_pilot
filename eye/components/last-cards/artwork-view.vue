<template>
  <div class="image-container">
    <a :href="`/card/${artwork.ID}`" @click.prevent="select()">
      <img :src="'/api/image/m/' + artwork.ID"/>
    </a>
    <div class="version">
      <span class="tag is-primary is-light">{{ artwork.Version }}</span>
    </div>
    <p class="is-size-7 tags">{{ artwork.Spell.Tags }}</p>
    <p class="is-size-7 info">id={{ artwork.ID }}, <b>seed: {{ artwork.Spell.Seed }}</b></p>
  </div>
</template>
<script>
export default {
  props: ['artwork'],
  methods: {
    select() {
      this.$emit('select')
    }
  }
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
    opacity: 0.7;
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
  color: #d4d1c3;
  visibility: hidden;
}

.image-container:hover p.tags {
  visibility: visible;
}

.image-container p.info {
  top: 30px;
  width: 100%;
  text-align: center;
  color: #d4d1c3;
  visibility: hidden;
}

.image-container:hover p.info {
  visibility: visible;
}

</style>
