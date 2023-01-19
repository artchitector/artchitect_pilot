<template>
  <section class="content">
    <top/>
    <philosophy/>
    <technics/>
    <launch/>
    <pray/>
    <next/>
  </section>
</template>
<script>
import Top from "@/components/idea/top.vue";
import Philosophy from "@/components/idea/philosophy.vue";
import Technics from "@/components/idea/technics.vue";
import Launch from "@/components/idea/launch.vue";
import Pray from "@/components/idea/pray.vue";
import Next from "@/components/idea/next.vue";

export default {
  components: {Next, Pray, Launch, Technics, Philosophy, Top}
}
</script>
