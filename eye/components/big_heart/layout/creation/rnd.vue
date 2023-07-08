<template>
  <div class="rnd">
    <NuxtLink :to="localePath(`/dream/${cardId}`)">
      <img :src="`/api/image/f/${cardId}`"/>
    </NuxtLink>
    <NuxtLink :to="localePath(`/dream/${cardId}`)" class="has-text-info">
      #{{ cardId }}
    </NuxtLink>
  </div>
</template>

<script>
export default {
  name: "rnd",
  props: ["cardId"]
}
</script>

<style lang="scss" scoped>
.rnd {
  margin-bottom: 30px;
}
</style>
