<template>
  <section class="content">
    <top/>
    <philosophy/>
    <technics/>
    <launch/>
    <prayer/>
    <next/>
  </section>
</template>
<script>
import Top from "@/components/idea/top.vue";
import Philosophy from "@/components/idea/philosophy.vue";
import Technics from "@/components/idea/technics.vue";
import Launch from "@/components/idea/launch.vue";
import Prayer from "@/components/idea/prayer.vue";
import Next from "@/components/idea/next.vue";

export default {
  components: {Next, Prayer, Launch, Technics, Philosophy, Top}
}
</script>
