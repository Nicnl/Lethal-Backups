<template>
  <div class="modal is-active">
    <div class="modal-background"  @click="close"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">
          {{ entry.infos.nbLoots + ' items ◽' + entry.infos.totalLootValue }}
        </p>
        <button class="delete" aria-label="close" @click="close"></button>
      </header>
      <section class="modal-card-body">

        <div class="item-icon" v-for="(amount, equipmentName) in extractItems(entry)" style="height: 80px;margin-bottom: 16px;">
          <img
              :src="$axios.defaults.baseURL + '/item_icon/' + equipmentName + '.webp'"
              style="height: 70px;padding: 8px;"
          >
          <span class="is-family-monospace item-amount">x{{ amount }}</span><br/>
        </div>

      </section>
      <footer class="modal-card-foot">
        <div class="buttons">
          <button class="button" @click="close">Close</button>
        </div>
      </footer>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.item-icon {
  display: inline-block;
  margin-right: 16px;
  background-color: rgba(100, 255, 100, 0.05);
  border-radius: 6px;
  position: relative;

  .item-amount {
    display: block;
    position: absolute;
    bottom: 0;
    right: 5px;
    font-size: 20px;

    opacity: 1.0;

    font-weight: bold;
  }
}
</style>

<script>
export default {
  name: 'ModalListItems',

  props: ['indexedItems', 'entry'],
  emits: ['close'],

  methods: {
    close() {
      this.$emit('close');
    },

    extractItems(entry) {
      let equipment = {};

      // entry.infos.shipGrabbableItemIDs.value is an array
      for (let id of entry.infos.shipGrabbableItemIDs.value) {
        let itemName = 'Unknown (ID: ' + id + ')';

        if (id in this.indexedItems) {
          itemName = this.indexedItems[id].name;

          if (this.indexedItems[id].tool) continue;
        } else {
          // May not be an equipment, skip
          continue;
        }

        if (itemName in equipment)
          equipment[itemName]++;
        else
          equipment[itemName] = 1;
      }

      return equipment;
    },
  },
}
</script>