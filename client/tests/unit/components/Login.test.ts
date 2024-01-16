import type { Mock } from 'vitest';
import { render, screen } from '@testing-library/vue';
import { RouterLinkStub } from '@vue/test-utils';
import { useRouter } from 'vue-router';
import { createTestingPinia } from '@pinia/testing';
import userEvent from '@testing-library/user-event';
import axios from 'axios';

import postLogin from '@/api/user';
import Login from '@/pages/Login.vue';

vi.mock('vue-router');
const useRouterMock = useRouter as Mock;

vi.mock('axios');
const axiosPostMock = axios.post as Mock;

const pinia = createTestingPinia();

describe('Login', () => {
  describe('when user submits login form', () => {
    it('directs user to home page upon successful login', async () => {
      const push = vi.fn();
      useRouterMock.mockReturnValue({ push });

      axiosPostMock.mockResolvedValue({
        data: {
          username: 'jiewei',
          email: 'jiewei@gmail.com',
        },
        status: 200,
      });

      render(Login, {
        global: {
          plugins: [pinia],
          stubs: {
            RouterLink: RouterLinkStub,
          },
        },
      });

      const emailInput = screen.getByRole('textbox', {
        name: /email/i,
      });
      await userEvent.type(emailInput, 'jiewei@email.com');

      const passwordInput = screen.getByLabelText(/password/i);
      await userEvent.type(passwordInput, 'Supersecretpassword000@');

      const submitButton = screen.getByRole('button', {
        name: /login/i,
      });
      await userEvent.click(submitButton);

      const body = {
        email: 'jiewei@email.com',
        password: 'Supersecretpassword000@',
      };
      expect(axios.post).toHaveBeenCalledWith(
        `http://unit-test/login`,
        JSON.stringify(body),
        { withCredentials: true }
      );

      expect(push).toHaveBeenCalledWith({ name: 'Home' });
    });

    it('login with unexpected error', async () => {
      axiosPostMock.mockRejectedValueOnce({
        data: {
          message: 'Unexpected error occurred',
        },
        status: 400,
      });

      render(Login, {
        global: {
          plugins: [pinia],
          stubs: {
            RouterLink: RouterLinkStub,
          },
        },
      });

      const emailInput = screen.getByRole('textbox', { name: /email/i });
      await userEvent.type(emailInput, 'jiewei@email.com');

      const passwordInput = screen.getByLabelText(/password/i);
      await userEvent.type(passwordInput, 'Supersecretpassword000@');

      const submitButton = screen.getByRole('button', { name: /login/i });
      await userEvent.click(submitButton);

      // Assert that the error message is displayed to the user
      expect(
        screen.getByText(/Unexpected error occurred/i)
      ).toBeInTheDocument();
    });
  });
});
