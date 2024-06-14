import type { RawAxiosResponseHeaders, AxiosResponseHeaders } from "axios";

export interface ApiResponse<T> {
  data: T;
  headers?: RawAxiosResponseHeaders | AxiosResponseHeaders;
  status: number;
}

export interface HttpStatusResponse {
  status: number;
}

export interface Pagination {
  page: number;
  pageSize: number;
  totalRecords?: number;
  totalPages?: number;
  hasNextPage?: boolean;
  hasPreviousPage?: boolean;
}
