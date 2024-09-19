<template>
  <div class="modal is-active">
    <div class="modal-background" @click="close"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">
          Restore backup
        </p>
        <button v-if="!isRestoring" class="delete" aria-label="close" @click="close"></button>
      </header>
      <section class="modal-card-body">

        <div class="content">
          <ul style="font-size: 18px;">
            <li>
              <b style="margin-right: 10px;">Planet:</b>
              <template v-if="entry.infos.CurrentPlanetID.value + '' in planets">
                {{ planets[entry.infos.CurrentPlanetID.value + ''] }}
              </template>
              <template>
                <span class="has-text-danger">Unknown planet {{ '(' + entry.infos.CurrentPlanetID.value + ')'}}</span>
              </template>
            </li>

            <li>
              <b style="margin-right: 5px;">Money:</b>
              {{ '◽' + entry.infos.GroupCredits.value }}
            </li>

            <li>
              <b style="margin-right: 5px;">Day:</b>
              {{ entry.infos.Stats_DaysSpent.value }}
            </li>

            <li>
              <b style="margin-right: 5px;">Deadline:</b>
              {{ Math.floor(entry.infos.DeadlineTime.value/0.75/60/24) + ' Days' }}
            </li>

            <li>
              <b style="margin-right: 5px;">Quota:</b>
              {{ '◽' + entry.infos.QuotaFulfilled.value + ' / ' + '◽' + entry.infos.ProfitQuota.value }}
            </li>

            <li>
              <b style="margin-right: 5px;">Scrap:</b>
              {{ entry.infos.nbLoots + ' items ◽' + entry.infos.totalLootValue }}
            </li>

            <li>
              <b style="margin-right: 5px;">Equipment:</b>
              <ul>
                <li v-for="(amount, equipmentName) in extractEquipment(entry)">
                  <div class="item-icon" style="height: 60px;margin-bottom: 16px;">
                    <img
                        :src="$axios.defaults.baseURL + '/item_icon/' + equipmentName + '.webp'"
                        style="height: 60px;padding: 5px;"
                    >
                    <span class="is-family-monospace item-amount">x{{ amount }}</span><br/>
                  </div>
                </li>
              </ul>
            </li>

          </ul>
        </div>

      </section>
      <footer class="modal-card-foot" style="justify-content: space-between;">
          <button class="button" @click="close" :disabled="isRestoring">Cancel</button>
          <button class="button is-success" @click="restoreBackup" :disabled="isRestoring" :class="{'is-loading': isRestoring}">
            <i class="fas fa-undo" style="margin-right: 5px;"></i>
            Restore
          </button>
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
  name: 'ModalRestore',

  props: ['indexedItems', 'planets', 'entry'],
  emits: ['close'],

  data() {
    return {
      isRestoring: false,
    };
  },

  methods: {
    close() {
      if (this.isRestoring) return;
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

    extractEquipment(entry) {
      let equipment = {};

      // entry.infos.shipGrabbableItemIDs.value is an array
      for (let id of entry.infos.shipGrabbableItemIDs.value) {
        let itemName = 'Unknown (ID: ' + id + ')';

        if (id in this.indexedItems) {
          itemName = this.indexedItems[id].name;

          if (!this.indexedItems[id].tool) continue;
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

    async restoreBackup() {
      if (this.isRestoring) return;
      this.isRestoring = true;

      try {
        await this.$axios.post('/restore/' + this.entry.hash);
        this.$notify({type: 'success', text: 'Backup restored successfully.'});
        this.isRestoring = false;
        this.close();
      } catch (error) {
        console.error(error);
        this.$notify({type: 'error', text: 'Error while restoring backup.'});
      } finally {
        this.isRestoring = false;
      }
    }
  },
}
</script>
