<template>
  <div>
    <section>
      <prayer/>
    </section>
  </div>
</template>
<script>
import Prayer from "@/components/pray/prayer";

export default {
  components: {Prayer},
  head() {
    return {
      title: this.$t('pray_title')
    }
  },
}
</script>
