import type { RawAxiosResponseHeaders, AxiosResponseHeaders } from "axios";

export interface APIResponse<T> {
  data: T;
  headers?: RawAxiosResponseHeaders | AxiosResponseHeaders;
  status: number;
}

export interface GENERIC_STATUS_RESPONSE {
  status: number;
}

export interface PAGINATION {
  page: number;
  pageSize: number;
  totalRecords?: number;
  totalPages?: number;
  hasNextPage?: boolean;
  hasPreviousPage?: boolean;
}
