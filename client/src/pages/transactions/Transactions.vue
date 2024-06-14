<template>
  <div class="relative">
    <div class="mb-4 flex items-center justify-between">
      <h1 class="text-xl font-bold tracking-wider">Transactions</h1>
      <div class="flex">
        <button
          class="flex items-center gap-1 rounded-2xl p-2 text-xs uppercase transition hover:bg-slate-100 dark:hover:bg-slate-700"
          :disabled="!pagination.hasPreviousPage"
          @click="handlePageNumberClick(1)"
          :class="
            !pagination.hasPreviousPage
              ? 'cursor-not-allowed text-gray-400 dark:text-gray-600'
              : ''
          "
        >
          <svg-icon type="mdi" class="h-5" :path="mdiArrowCollapseLeft" />
        </button>
        <button
          class="mr-4 flex items-center gap-1 rounded-2xl p-2 text-xs uppercase transition hover:bg-slate-100 dark:hover:bg-slate-700"
          :disabled="!pagination.hasPreviousPage"
          @click="handlePageNumberClick(pagination.page - 1)"
          :class="
            !pagination.hasPreviousPage
              ? 'cursor-not-allowed text-gray-400 dark:text-gray-600'
              : ''
          "
        >
          <svg-icon type="mdi" :path="mdiChevronLeft" />
        </button>
        <div class="flex gap-1">
          <button
            v-for="pageNumber in computedPageNumbers"
            :key="pageNumber?.toString()"
            class="rounded-lg px-4 py-2 text-sm transition"
            :class="dynamicPaginationClassNames(pageNumber)"
            @click="handlePageNumberClick(pageNumber)"
          >
            {{ pageNumber }}
          </button>
        </div>
        <button
          class="ml-4 flex items-center gap-1 rounded-2xl p-2 text-xs uppercase hover:bg-gray-100 dark:hover:bg-slate-700"
          :disabled="!pagination.hasNextPage"
          :class="
            !pagination.hasNextPage
              ? 'cursor-not-allowed text-gray-400 dark:text-gray-600'
              : ''
          "
          @click="handlePageNumberClick(pagination.page + 1)"
        >
          <svg-icon type="mdi" :path="mdiChevronRight" />
        </button>
        <button
          class="flex items-center gap-1 rounded-2xl p-2 text-xs uppercase hover:bg-gray-100 dark:hover:bg-slate-700"
          :disabled="!pagination.hasNextPage"
          :class="
            !pagination.hasNextPage
              ? 'cursor-not-allowed text-gray-400 dark:text-gray-600'
              : ''
          "
          @click="handlePageNumberClick(pagination.totalPages!)"
        >
          <svg-icon class="h-5" type="mdi" :path="mdiArrowCollapseRight" />
        </button>
      </div>
    </div>

    <div class="mb-4 overflow-x-auto">
      <table
        class="border border-gray-200 bg-white dark:border-gray-500 dark:bg-gray-800"
      >
        <thead>
          <tr class="bg-gray-200 text-sm leading-normal dark:bg-gray-700">
            <th class="px-6 py-3 text-center">Sender</th>
            <th class="px-6 py-3 text-center">Beneficiary</th>
            <th class="px-6 py-3 text-center">Source Amount</th>
            <th class="px-6 py-3 text-center">Destination Amount</th>
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
              {{ transaction.sender_username }}&nbsp;({{
                transaction.sender_mobile_number
              }})
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.beneficiary_username }}&nbsp;({{
                transaction.beneficiary_mobile_number
              }})
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.source_amount }}
              {{ transaction.source_currency }}
            </td>
            <td class="px-6 py-3 text-center">
              {{ transaction.destination_amount }}
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
      <div class="text-sm text-slate-700 dark:text-slate-300">
        Showing {{ showingRecordsComputation }} of
        {{ pagination.totalRecords }} records.
      </div>

      <div class="flex">
        <button
          class="flex items-center gap-1 rounded-2xl p-2 text-xs uppercase transition hover:bg-slate-100 dark:hover:bg-slate-700"
          :disabled="!pagination.hasPreviousPage"
          @click="handlePageNumberClick(1)"
          :class="
            !pagination.hasPreviousPage
              ? 'cursor-not-allowed text-gray-400 dark:text-gray-600'
              : ''
          "
        >
          <svg-icon type="mdi" class="h-5" :path="mdiArrowCollapseLeft" />
        </button>
        <button
          class="mr-4 flex items-center gap-1 rounded-2xl p-2 text-xs uppercase transition hover:bg-slate-100 dark:hover:bg-slate-700"
          :disabled="!pagination.hasPreviousPage"
          @click="handlePageNumberClick(pagination.page - 1)"
          :class="
            !pagination.hasPreviousPage
              ? 'cursor-not-allowed text-gray-400 dark:text-gray-600'
              : ''
          "
        >
          <svg-icon type="mdi" :path="mdiChevronLeft" />
        </button>
        <div class="flex gap-1">
          <button
            v-for="pageNumber in computedPageNumbers"
            :key="pageNumber?.toString()"
            class="rounded-lg px-4 py-2 text-sm transition"
            :class="dynamicPaginationClassNames(pageNumber)"
            @click="handlePageNumberClick(pageNumber)"
          >
            {{ pageNumber }}
          </button>
        </div>
        <button
          class="ml-4 flex items-center gap-1 rounded-2xl p-2 text-xs uppercase hover:bg-gray-100 dark:hover:bg-slate-700"
          :disabled="!pagination.hasNextPage"
          :class="
            !pagination.hasNextPage
              ? 'cursor-not-allowed text-gray-400 dark:text-gray-600'
              : ''
          "
          @click="handlePageNumberClick(pagination.page + 1)"
        >
          <svg-icon type="mdi" :path="mdiChevronRight" />
        </button>
        <button
          class="flex items-center gap-1 rounded-2xl p-2 text-xs uppercase hover:bg-gray-100 dark:hover:bg-slate-700"
          :disabled="!pagination.hasNextPage"
          :class="
            !pagination.hasNextPage
              ? 'cursor-not-allowed text-gray-400 dark:text-gray-600'
              : ''
          "
          @click="handlePageNumberClick(pagination.totalPages!)"
        >
          <svg-icon class="h-5" type="mdi" :path="mdiArrowCollapseRight" />
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
import type { Pagination } from "@/types/generic";
import {
  mdiArrowLeft,
  mdiArrowRight,
  mdiArrowCollapseLeft,
  mdiArrowCollapseRight,
  mdiChevronLeft,
  mdiChevronRight,
} from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";

const toastStore = useToastStore();
const transactions = ref<Transaction[]>([]);
const pagination = ref<Pagination>({
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

const showingRecordsComputation = computed(() => {
  let recordsComputation = pagination.value.pageSize * pagination.value.page;

  if (recordsComputation > pagination.value.totalRecords!) {
    return pagination.value.totalRecords;
  }
  return recordsComputation;
});

const computedPageNumbers = computed(() => {
  const currentPageNumber = pagination.value.page;
  const totalPages = pagination.value.totalPages;

  let start = currentPageNumber - 2;
  let end = currentPageNumber + 2;

  if (start < 1) {
    start = 1;
    end = start + 4;
  }

  if (end > totalPages!) {
    end = totalPages!;
    start = end - 4;
    if (start < 1) {
      start = 1;
    }
  }

  const pageNumbers: number[] = [];
  for (let i = start; i <= end; i++) {
    pageNumbers.push(i);
  }

  return pageNumbers;
});

const dynamicPaginationClassNames = (pageNumber: number) => {
  if (pageNumber === 0) {
    return "text-black dark:text-white";
  }

  if (pageNumber !== pagination.value.page) {
    return "text-black hover:bg-slate-200 dark:text-white dark:hover:bg-slate-700";
  }

  return "bg-slate-900 text-white hover:bg-slate-700 dark:bg-slate-200 dark:text-slate-700";
};

const handlePageNumberClick = async (pageNumber: number) => {
  if (pageNumber === 0) {
    return;
  }

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
    case "success":
      return "text-green-500 dark:text-green-400";
    case "failed":
      return "text-red-500 dark:text-red-400";
    default:
      return "text-orange-500 text-orange-400";
  }
};

const getTransactions = async () => {
  try {
    // retrieve url params
    const queryParams: Pagination = {
      page: parseInt(route.query?.page as string) || 1,
      pageSize: parseInt(route.query?.pageSize as string) || 10,
    };

    const { data, headers, status } = await GET_TRANSACTIONS(queryParams);

    if (status === 200) {
      transactions.value = data;

      if (headers) {
        pagination.value = {
          page: parseInt(headers["x-page"]) || 1,
          pageSize: parseInt(headers["x-page-size"]) || 1,
          totalRecords: parseInt(headers["x-total"]) || 1,
          totalPages: parseInt(headers["x-total-pages"]) || 1,
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
