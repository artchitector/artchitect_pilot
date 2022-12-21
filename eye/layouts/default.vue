<template>
  <div>
    <div class="has-text-centered is-header pt-1 pb-2">
      <NuxtLink :to="{name: 'index'}">
        <h1 class="is-size-3 has-text-weight-bold">
          artchitect.space
        </h1>
      </NuxtLink>
    </div>

    <div class="container is-max-desktop">
      <Nuxt/>
    </div>

    <footer class="footer">
      <p class="has-text-centered">
        github: <a href="https://github.com/artchitector/artchitect" target="_blank">https://github.com/artchitector/artchitect</a>
      </p>
    </footer>
  </div>
</template>
<script>

</script>
<style lang="scss">
.is-header {
  background-color: #2d2d2d;
  h1 {
    color: #d4d1c3;
  }
}
footer.footer {
  background-color: #2d2d2d;
  color: hsl(0, 0%, 96%);
}
</style>
