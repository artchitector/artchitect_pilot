<template>
  <div>
    <a href="#" @click.prevent="jsonVisible = !jsonVisible">
      <template v-if="jsonVisible">
        hide json
      </template>
      <template v-else>
        show json
      </template>
    </a>
    <pre v-if="jsonVisible" class="mb-2">{{ state }}</pre>
    <table class="table is-fullwidth">
      <tr>
        <td>Current GAN state</td>
        <td>{{ state.CurrentState }}</td>
      </tr>
      <tr>
        <td></td>
        <td></td>
      </tr>
      <tr>
        <th colspan="2" class="has-text-centered is-selected">
          Last painting
        </th>
      </tr>
      <tr>
        <td colspan="2" class="has-text-centered">
          <template v-if="!state.LastPainting">
            no painting yet...
          </template>
          <img v-else :src="'/api/painting/' + state.LastPainting.ID"/>
        </td>
      </tr>
    </table>
  </div>
</template>

<script>
export default {
  props: {
    // eslint-disable-next-line vue/require-default-prop
    state: Object
  },
  data: () => ({
    jsonVisible: false
  })
}
</script>
