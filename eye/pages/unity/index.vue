<i18n>
{
  "en": {
    "title": "Artchitect - unity",
    "subtitle": "unity"
  },
  "ru": {
    "title": "Artchitect - единство",
    "subtitle": "единство"
  }
}
</i18n>

<template>
  <section>
    <h3 class="has-text-centered is-size-4">{{ $t('subtitle') }}</h3>
    <template v-if="$fetchState.pending">
      <div class="has-text-centered">
        <loader/>
      </div>
    </template>
    <template v-else-if="$fetchState.error">
      <div class="notification is-danger">
        {{ $fetchState.error.message }}
      </div>
    </template>
    <template v-else>
      <unity-list :unities="unities" visible-count="50" cards-in-column="3"/>
    </template>
  </section>
</template>

<script>
import UnityList from "@/components/unity/unity-list.vue";

export default {
  name: "unity",
  components: {UnityList},
  head() {
    return {
      title: this.$t('title')
    }
  },
  data() {
    return {
      unities: []
    };
  },
  async fetch() {
    this.unities = await this.$axios.$get("/unity")
  }
}
</script>

<style scoped>

</style>
