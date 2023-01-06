<template>
  <div class="image-container">
    <a :href="'/api/painting/' + artwork.ID" target="_blank">
      <img :src="'/api/painting/' + artwork.ID"/>
      <p class="is-size-7 tags">{{artwork.Spell.Tags}}</p>
      <p class="is-size-7 info">id={{artwork.ID}}, <b>seed: {{artwork.Spell.Seed}}</b></p>
    </a>
  </div>
</template>
<script>
export default {
  props: ['artwork']
}
</script>
<style>
.image-container {
  position: relative;
}
.image-container p {
  visibility: hidden;
  position: absolute;
  margin: auto;
  word-break: break-all;
  background-color: rgba(0, 0, 0, 0.6);
  padding: 5px;
}
.image-container:hover p {
  visibility: visible;
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
