import { API_ENDPOINTS } from "./endpoints";
import { apiClient } from "./runtime";
import type {
  CommunityAccount,
  CommunityAccountPayload,
  CommunityReport,
  CommunityReportPayload,
  CommunityReviewReportInput,
  CommunityReviewSubmissionInput,
  CommunitySubmission,
  CommunitySubmissionPayload,
  CommunityUpdateAccountInput,
} from "./types";

type RequestOptions = {
  signal?: AbortSignal;
};

export type CommunityAccountListQuery = {
  keyword?: string;
  limit?: number;
  role?: string;
  status?: string;
};

export type CommunityReviewQueueQuery = {
  limit?: number;
  status?: string;
};

export const communityApi = {
  listAccounts: (query: CommunityAccountListQuery = {}, options: RequestOptions = {}) =>
    apiClient.request<CommunityAccountPayload>(API_ENDPOINTS.community.accounts, {
      query,
      signal: options.signal,
    }),
  listReports: (query: CommunityReviewQueueQuery = {}, options: RequestOptions = {}) =>
    apiClient.request<CommunityReportPayload>(API_ENDPOINTS.community.reports, {
      query,
      signal: options.signal,
    }),
  listSubmissions: (query: CommunityReviewQueueQuery = {}, options: RequestOptions = {}) =>
    apiClient.request<CommunitySubmissionPayload>(API_ENDPOINTS.community.submissions, {
      query,
      signal: options.signal,
    }),
  reviewReport: (
    reportId: number | string,
    body: CommunityReviewReportInput,
    options: RequestOptions = {},
  ) =>
    apiClient.request<CommunityReport>(API_ENDPOINTS.community.report(reportId), {
      body,
      method: "PATCH",
      signal: options.signal,
    }),
  reviewSubmission: (
    submissionId: number | string,
    body: CommunityReviewSubmissionInput,
    options: RequestOptions = {},
  ) =>
    apiClient.request<CommunitySubmission>(API_ENDPOINTS.community.submissionReview(submissionId), {
      body,
      method: "PATCH",
      signal: options.signal,
    }),
  updateAccount: (
    accountId: number | string,
    body: CommunityUpdateAccountInput,
    options: RequestOptions = {},
  ) =>
    apiClient.request<CommunityAccount>(API_ENDPOINTS.community.account(accountId), {
      body,
      method: "PATCH",
      signal: options.signal,
    }),
};
