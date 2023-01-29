<template>
    <img class="loader-image" :class="`loader-image-size-${size}`" src="/loader.png"/>
</template>
<script>
export default {
  props: ['size']
}
</script>
<style>
.loader-image {
  animation: spinAround 1300ms infinite linear;
}

.loader-image-size-undefined {
  height: 120px;
  width: 120px;
}

.loader-image-size-s {
  height: 2rem;
  width: 2rem;
}
</style>
