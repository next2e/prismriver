<template>
  <li>
    <span class="text-overflow">
      <button v-on:click="queue" class="select-song">Add</button>
      {{ title }}
    </span>
  </li>
</template>

<script lang="ts">
  import { Component, Mixins, Prop } from 'vue-property-decorator'

  import ShowMessage from '../mixins/ShowMessage'

  @Component
  export default class SearchItem extends Mixins(ShowMessage) {
    @Prop(String) id!: string
    @Prop(String) title!: string
    @Prop(String) type!: string

    queue (): void {
      this.showMessage('Adding ' + this.title)
      const params = new URLSearchParams()
      params.set('id', this.id)
      params.set('type', this.type)
      fetch('queue', {
        body: params,
        method: 'POST'
      })
    }
  }
</script>