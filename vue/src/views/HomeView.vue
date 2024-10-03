<template>
  <main>
    <div v-if="isLoading" class="has-text-centered" style="padding-top: 32px;">
      Loading backups...<br/><br/>
      <button class="button is-loading">---</button>
    </div>

    <template v-if="backupSlots !== null">
      <div class="tabs is-large is-centered" style="margin-bottom: 16px;">
        <ul>
          <li
              v-for="slot in slots"
              :class="{'is-active': activeTab === slot.slot}"
          ><a @click="activeTab = slot.slot">{{ slot.name }}</a></li>
        </ul>
      </div>

      <table class="table is-fullwidth is-striped is-hoverable">
        <thead>
        <tr>
          <th style="width: 20px;"></th>
          <th style="width: 200px;">Date</th>
          <th style="width: 180px;">Planet</th>
          <th style="width: 100px;">Money</th>
          <th style="width: 85px;">Day</th>
          <th style="width: 95px;">Deadline</th>
          <th style="width: 110px;">Quota</th>
          <th style="width: 150px;">Scrap</th>
          <th>Equipment</th>
        </tr>
        </thead>
        <tbody>

        <tr v-for="entry in backupSlots[activeTab]" style="height: 89px;">

          <th class="is-vcentered">
            <button class="button is-success is-outlined" @click="showModalRestore = entry" style="width: 35px;">
              <i class="fas fa-undo"></i>
            </button>
          </th>

          <th class="is-vcentered is-family-monospace">{{ entry.time.replaceAll('T', ' ').replaceAll('Z', '') }}</th>
          <th class="is-vcentered">
            <template v-if="entry.infos.CurrentPlanetID.value + '' in planets">
              {{ planets[entry.infos.CurrentPlanetID.value + ''] }}
            </template>
            <template>
              <span class="has-text-danger">Unknown planet {{ '(' + entry.infos.CurrentPlanetID.value + ')'}}</span>
            </template>
          </th>
          <th class="is-vcentered is-family-monospace has-text-primary">{{ '■ ' + entry.infos.GroupCredits.value }}</th>
          <th class="is-vcentered is-family-monospace">{{ 'Day ' + entry.infos.Stats_DaysSpent.value }}</th>
          <th class="is-vcentered is-family-monospace has-text-warning">{{ Math.floor(entry.infos.DeadlineTime.value/0.75/60/24) + ' Days' }}</th>
          <th class="is-vcentered is-family-monospace has-text-primary">
            <template v-if="entry.infos.QuotaFulfilled.value > 0">
              {{ '■ ' + entry.infos.QuotaFulfilled.value + ' / ' + '■ ' + entry.infos.ProfitQuota.value }}
            </template>
            <template v-else>
              {{ '■ ' + entry.infos.ProfitQuota.value }}
            </template>
          </th>

          <th class="is-vcentered">
            <button
                class="button is-info is-outlined"
                style="margin-right: 8px;"
                @click="showModalItems = entry"
            >
              {{ nbLoots(extractItems(entry)) + ' items' }}
              <span class="has-text-primary" style="margin-left: 8px;">
                {{ '■ ' + entry.infos.totalLootValue }}
              </span>
            </button>
          </th>

          <th class="is-vcentered">
            <div
                v-for="(amount, equipmentName) in extractEquipment(entry)"
                class="item-icon"
                style="height: 66px;"
            >
              <img
                  :src="$axios.defaults.baseURL + '/item_icon/' + equipmentName + '.webp'"
                  style="height: 66px;"
              >
              <span class="is-family-monospace item-amount">x{{ amount }}</span><br/>
            </div>
          </th>
        </tr>
        </tbody>
      </table>
    </template>
  </main>

  <ModalListItems
      v-if="showModalItems !== null"
      :indexed-items="indexedItems"
      :items="extractItems(showModalItems)"
      :entry="showModalItems"
      @close="showModalItems = null"
  />

  <ModalRestore
      v-if="showModalRestore !== null"
      :indexed-items="indexedItems"
      :planets="planets"
      :entry="showModalRestore"
      @close="showModalRestore = null"
  />
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
    font-size: 18px;

    opacity: 0.8;
  }
}
</style>

<script>
import ModalListItems from "@/components/ModalListItems.vue";
import ModalRestore from "@/components/ModalRestore.vue";

export default {
  name: 'HomeView',

  data() {
    return {
      isLoading: false,

      backupSlots: null,
      planets: null,
      items: null,

      activeTab: "slot1",
      showModalItems: null,
      showModalRestore: null,
    };
  },

  async created() {
    this.isLoading = true;

    await this.loadPlanets();
    await this.loadItems();
    await this.loadSaves();
  },

  computed: {
    slots() {
      if (this.backupSlots === null) return null;

      let output = [];


      let slots = Object.keys(this.backupSlots);
      slots.sort()

      for (let slot of slots) {
        output.push({
          slot: slot,
          name: slot.replaceAll('slot', 'Save '),
        });
      }


      return output;
    },

    indexedItems() {
      let output = {};

      for (let item of this.items) output[item.id] = item;

      return output;
    }
  },

  methods: {
    async loadPlanets() {
      try {
        const resp = await this.$axios.get('/planets');
        this.planets = resp.data;
      } catch (error) {
        console.error(error);
        this.$notify({type: 'error', text: 'Error when listing planets.'});
      } finally {
        this.isLoading = false;
      }
    },

    async loadItems() {
      try {
        const resp = await this.$axios.get('/items');
        this.items = resp.data;
      } catch (error) {
        console.error(error);
        this.$notify({type: 'error', text: 'Error when listing items.'});
      } finally {
        this.isLoading = false;
      }
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

    extractItems(entry) {
      let equipment = {};

      // entry.infos.shipGrabbableItemIDs.value is an array
      for (let id of entry.infos.shipGrabbableItemIDs.value) {
        let itemName = 'unknown';

        if (id in this.indexedItems) {
          if (this.indexedItems[id].tool) continue;
          itemName = this.indexedItems[id].name;
        }

        if (itemName in equipment)
          equipment[itemName]++;
        else
          equipment[itemName] = 1;
      }

      return equipment;
    },

    nbLoots(items) {
      let nb = 0;
      for (const name in items)
        nb += items[name];
      return nb;
    },

    async loadSaves() {
      if (this.showModalItems !== null || this.showModalRestore !== null) {
        setTimeout(() => {
          this.loadSaves();
        }, 5000);
        return;
      }

      try {
        const resp = await this.$axios.get('/saves');

        let backupSlots = {}; // {slot1: [backups], slot2: [backups], ...}
        for (let k in resp.data) {
          let backup = resp.data[k];
          if (!backupSlots[backup.slot]) backupSlots[backup.slot] = [];

          let nbLoots = 0;
          let totalLootValue = 0;

          if ('shipGrabbableItemIDs' in backup.infos)
            nbLoots = backup.infos.shipGrabbableItemIDs.value.length;

          if ('shipScrapValues' in backup.infos)
            totalLootValue = backup.infos.shipScrapValues.value.reduce((a, b) => a + b, 0);

          backup.infos.nbLoots = nbLoots;
          backup.infos.totalLootValue = totalLootValue;

          backupSlots[backup.slot].push(backup);
        }

        // Sort using "unixTime" field
        for (let slot in backupSlots)
          backupSlots[slot].sort((a, b) => b.unixTime - a.unixTime);

        this.backupSlots = backupSlots;
        this.isLoading = false;
      } catch (error) {
        this.isLoading = true;
        this.backupSlots = null;
        console.error(error);
        this.$notify({type: 'error', text: 'Error when listing backups.'});
      } finally {
        setTimeout(() => {
          this.loadSaves();
        }, 5000);
      }
    }
  },

  components: {
    ModalListItems,
    ModalRestore,
  }
}
</script>
