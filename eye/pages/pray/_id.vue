<template>
  <section>
    <answer :id="this.$route.params.id"/>
  </section>
</template>
<script>
import Answer from "@/components/pray/answer.vue";

export default {
  components: {Answer},
  head () {
    return {
      title: `Artchitect - pray answer #${this.$route.params.id}`
    }
  }
}
</script>
