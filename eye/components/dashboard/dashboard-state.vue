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
        <td>Current Artchitect state</td>
        <td>
          {{ state.CurrentState.State }}
        </td>
      </tr>
      <tr>
        <td>Last decision made</td>
        <td>
          <template v-if="!!state.LastDecision">
            <span
              class="is-family-monospace is-size-7 has-text-grey">id: {{ state.LastDecision.ID }} ({{ state.LastDecision.CreatedAt }})</span><br/>
            {{ state.LastDecision.Result }}
            <img class="seedImage" :src="`data:image/jpeg;base64, ${state.LastDecision.Image}`"/>
          </template>
          <template v-else>
            no decision
          </template>
        </td>
      </tr>
      <tr>
        <td>Last spell (artwork keywords)</td>
        <td>
          <template v-if="!!state.LastSpell">
            <span
              class="is-family-monospace is-size-7 has-text-grey">id: {{ state.LastSpell.ID }} ({{ state.LastSpell.CreatedAt }})</span><br/>
            idea: {{ state.LastSpell.Idea }}<br/>
            tags for artist: {{ state.LastSpell.Tags }}<br/>
            seed: {{ state.LastSpell.Seed }}
          </template>
          <template v-else>
            no last spell
          </template>
        </td>
      </tr>
      <!-- empty line -->
      <tr>
        <td></td>
        <td></td>
      </tr>
      <tr>
        <th colspan="2" class="has-text-centered is-selected">
          Last painting (id={{ state.LastPainting.ID }}, spell_id={{ state.LastPainting.Spell.ID }})
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
