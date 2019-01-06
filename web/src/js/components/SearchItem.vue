<template>
  <li>
    <span class="text-overflow">
      <button v-on:click="queue" class="select-song">Add</button>
      {{ title }}
    </span>
  </li>
</template>

<script lang="ts">
  import $ from 'jquery'
  import { Component, Mixins, Prop } from 'vue-property-decorator'

  import ShowMessage from '../mixins/ShowMessage'

  @Component
  export default class SearchItem extends Mixins(ShowMessage) {
    @Prop(String) id!: string
    @Prop(String) title!: string
    @Prop(String) type!: string

    queue (): void {
      this.showMessage('Adding ' + this.title)
      $.post('/queue', { id: this.id, type: this.type })
    }
  }
</script>