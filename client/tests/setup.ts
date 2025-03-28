import { cleanup } from '@testing-library/vue';
import * as matchers from '@testing-library/jest-dom/matchers';
import { expect, afterEach } from 'vitest';

expect.extend(matchers);

afterEach(() => {
  // ensure each test starts with a clean state
  cleanup();
});
