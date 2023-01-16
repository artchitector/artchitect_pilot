<template>
  <section class="green-section has-text-centered">
    <p>
      🚀 <b>Artchitect</b> {{$t('launched')}}<br/>
      <a class="button is-success mt-2" href="#cards">{{ $t('to_cards')}}</a>
      <NuxtLink class="button is-success mt-2" :to="localePath('pray')">{{ $t('to_pray')}}</NuxtLink>
    </p>
  </section>
</template>
