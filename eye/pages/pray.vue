<template>
  <div>
    <section>
      <h1 class="is-size-3">{{$t('pray_place')}}</h1>
      <prayer/>
    </section>
  </div>
</template>
<script>
import Prayer from "@/components/prayer";

export default {
  components: {Prayer},
  head: {
    title: 'Artchitect - Молитва'
  },
}
</script>
