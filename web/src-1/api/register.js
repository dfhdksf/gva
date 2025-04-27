// filepath: src/api/register.js

import axios from 'axios';

const API_URL = 'your_api_endpoint_here'; // Replace with your actual API endpoint

// Function to send a verification code to the provided phone number
export const sendCode = async (phoneNumber) => {
  try {
    const response = await axios.post(`${API_URL}/send-verification-code`, {
      phoneNumber
    });
    return response.data;
  } catch (error) {
    throw new Error('Failed to send verification code');
  }
};