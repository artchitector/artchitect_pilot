<template>
  <section>
    <answer v-if="this.$route.params.id" :id="this.$route.params.id"/>
    <div v-else class="has-text-centered">
      <loader/>
    </div>
  </section>
</template>
<script>
import Answer from "@/components/pray/answer.vue";

export default {
  components: {Answer},
  head() {
    return {
      title: this.$route.params.id ?
        `Artchitect - pray answer #${this.$route.params.id}` :
        'Artchitect - pray to Universe'
    }
  },
  mounted() {
    if (!this.$route.params.id) {
      this.$router.push("/pray/create");
    }
  }
}
</script>
