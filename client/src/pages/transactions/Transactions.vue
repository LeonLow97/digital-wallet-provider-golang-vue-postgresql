<template>
  <div class="relative">
    <div class="mb-4 flex items-center justify-between">
      <h1 class="text-xl font-bold tracking-wider">Transactions</h1>
      <div class="flex">
        <button
          class="mr-4 flex items-center gap-1 rounded-2xl px-4 py-2 text-xs uppercase transition hover:bg-slate-100"
          :disabled="!pagination.hasPreviousPage"
          @click="handlePageNumberClick(pagination.page - 1)"
          :class="
            !pagination.hasPreviousPage
              ? 'cursor-not-allowed text-gray-400'
              : ''
          "
        >
          <svg-icon class="h-4" type="mdi" :path="mdiArrowLeft" /> Previous
        </button>
        <div class="flex gap-1">
          <button
            v-for="pageNumber in computedPageNumbers"
            :key="pageNumber?.toString()"
            class="rounded-lg px-4 py-2 text-sm transition"
            :class="
              pageNumber !== pagination.page
                ? 'text-black hover:bg-slate-200'
                : 'bg-slate-900 text-white hover:bg-slate-700'
            "
            @click="handlePageNumberClick(pageNumber)"
          >
            {{ pageNumber }}
          </button>
        </div>
        <button
          class="ml-4 flex items-center gap-1 rounded-2xl px-4 py-2 text-xs uppercase hover:bg-gray-100"
          :disabled="!pagination.hasNextPage"
          :class="
            !pagination.hasNextPage ? 'cursor-not-allowed text-gray-400' : ''
          "
          @click="handlePageNumberClick(pagination.page + 1)"
        >
          Next <svg-icon class="h-4" type="mdi" :path="mdiArrowRight" />
        </button>
      </div>
    </div>

    <div class="mb-4 overflow-x-auto">
      <table
        class="border border-gray-200 bg-white dark:border-gray-500 dark:bg-gray-800"
      >
        <thead>
          <tr class="bg-gray-200 text-sm leading-normal dark:bg-gray-700">
            <th class="px-6 py-3 text-center">Sender Username</th>
            <th class="px-6 py-3 text-center">Sender Mobile</th>
            <th class="px-6 py-3 text-center">Beneficiary Username</th>
            <th class="px-6 py-3 text-center">Beneficiary Mobile</th>
            <th class="px-6 py-3 text-center">Source Amount</th>
            <th class="px-6 py-3 text-center">Source Currency</th>
            <th class="px-6 py-3 text-center">Destination Amount</th>
            <th class="px-6 py-3 text-center">Destination Currency</th>
            <th class="px-6 py-3 text-center">Source of Transfer</th>
            <th class="px-6 py-3 text-center">Status</th>
            <th class="px-6 py-3 text-center">Created At</th>
          </tr>
        </thead>
        <tbody class="text-sm">
          <tr
            v-for="transaction in transactions"
            :key="transaction.created_at"
            class="border-b border-gray-200 hover:bg-gray-100 dark:border-gray-600 dark:hover:bg-gray-600"
          >
            <td class="px-6 py-3 text-center">
              {{ transaction.sender_username }}
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.sender_mobile_number }}
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.beneficiary_username }}
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.beneficiary_mobile_number }}
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.source_amount }}
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.source_currency }}
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.destination_amount }}
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.destination_currency }}
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.source_of_transfer }}
            </td>
            <td
              class="px-6 py-3 text-center"
              :class="transactionCssTextColor(transaction.status)"
            >
              {{ transaction.status }}
            </td>
            <td class="px-6 py-3 text-center">
              {{ new Date(transaction.created_at).toLocaleString() }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="mb-4 flex items-center justify-between">
      <div class="text-sm text-slate-700">
        Showing {{ pagination.pageSize }} of
        {{ pagination.totalRecords }} records.
      </div>
      <div class="flex">
        <button
          class="mr-4 flex items-center gap-1 rounded-2xl px-4 py-2 text-xs uppercase transition hover:bg-slate-100"
          :disabled="!pagination.hasPreviousPage"
          @click="handlePageNumberClick(pagination.page - 1)"
          :class="
            !pagination.hasPreviousPage
              ? 'cursor-not-allowed text-gray-400'
              : ''
          "
        >
          <svg-icon class="h-4" type="mdi" :path="mdiArrowLeft" /> Previous
        </button>
        <div class="flex gap-1">
          <button
            v-for="pageNumber in computedPageNumbers"
            :key="pageNumber?.toString()"
            class="rounded-lg px-4 py-2 text-sm transition"
            :class="
              pageNumber !== pagination.page
                ? 'text-black hover:bg-slate-200'
                : 'bg-slate-900 text-white hover:bg-slate-700'
            "
            @click="handlePageNumberClick(pageNumber)"
          >
            {{ pageNumber }}
          </button>
        </div>
        <button
          class="ml-4 flex items-center gap-1 rounded-2xl px-4 py-2 text-xs uppercase hover:bg-gray-100"
          :disabled="!pagination.hasNextPage"
          :class="
            !pagination.hasNextPage ? 'cursor-not-allowed text-gray-400' : ''
          "
          @click="handlePageNumberClick(pagination.page + 1)"
        >
          Next <svg-icon class="h-4" type="mdi" :path="mdiArrowRight" />
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { GET_TRANSACTIONS } from "@/api/transactions";
import type { Transaction } from "@/types/transactions";
import { computed, onMounted, ref } from "vue";
import { useToastStore } from "@/stores/toast";
import { useRoute, useRouter } from "vue-router";
import type { PAGINATION } from "@/types/generic";
import { mdiArrowLeft, mdiArrowRight } from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";

const toastStore = useToastStore();
const transactions = ref<Transaction[]>([]);
const pagination = ref<PAGINATION>({
  page: 1,
  pageSize: 10,
  totalRecords: 0,
  totalPages: 0,
  hasNextPage: false,
  hasPreviousPage: false,
});
const route = useRoute();
const router = useRouter();

onMounted(() => {
  // add page and page size to url query params when page loads
  if (!route.query.page || !route.query.pageSize) {
    router.replace({
      path: route.path,
      query: {
        page: route.query.page || 1,
        pageSize: route.query.pageSize || 10,
      },
    });
  }

  getTransactions();
});

const computedPageNumbers = computed(() => {
  const currentPageNumber = pagination.value.page;
  const totalPages = pagination.value.totalPages;

  let start = currentPageNumber - 2;
  let end = currentPageNumber + 2;

  if (start < 1) {
    start = 1;
  }

  const pageNumbers: number[] = [];
  for (let i = start; i < end + 1; i++) {
    if (totalPages && i <= totalPages) {
      pageNumbers.push(i);
    }
  }

  return pageNumbers;
});

const handlePageNumberClick = async (pageNumber: number) => {
  await router.replace({
    path: route.path,
    query: {
      page: pageNumber,
      pageSize: pagination.value.pageSize,
    },
  });

  // Call getTransactions to fetch data for the selected page
  getTransactions();
};

const transactionCssTextColor = (status: string) => {
  switch (status?.toLowerCase()) {
    case "successful":
      return "text-green-500";
    case "failed":
      return "text-red-500";
    default:
      return "text-blue-500";
  }
};

const getTransactions = async () => {
  try {
    // retrieve url params
    const queryParams: PAGINATION = {
      page: parseInt(route.query?.page as string) || 1,
      pageSize: parseInt(route.query?.pageSize as string) || 10,
    };

    const { data, headers, status } = await GET_TRANSACTIONS(queryParams);

    if (status === 200) {
      transactions.value = data;

      if (headers) {
        pagination.value = {
          page: parseInt(headers["x-page"]) || NaN,
          pageSize: parseInt(headers["x-page-size"]) || NaN,
          totalRecords: parseInt(headers["x-total"]) || NaN,
          totalPages: parseInt(headers["x-total-pages"]) || NaN,
          hasPreviousPage: headers["x-has-previous-page"] === "true",
          hasNextPage: headers["x-has-next-page"] === "true",
        };

        router.replace({
          path: route.path,
          query: {
            page: pagination.value.page,
            pageSize: pagination.value.pageSize,
          },
        });
      }
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};
</script>
