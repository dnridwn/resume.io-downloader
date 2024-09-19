<script setup>
    import axios from 'axios'
    import { ref } from 'vue'

    let isLoading = ref(false)

    const renderingToken = ref('')
    const errorMessage = ref('')

    function generate() {
        const renderingTokenInput = renderingToken.value;
        if (!renderingTokenInput || renderingTokenInput == '') {
            return;   
        }
        
        errorMessage.value = '';
        isLoading.value = true;

        axios({
            url: import.meta.env.VITE_APP_API_URL + '/generate',
            method: 'GET',
            responseType: 'blob',
            params: {
                rendering_token: renderingTokenInput
            }
        }).then((response) => {
            const href = URL.createObjectURL(response.data);
            const link = document.createElement('a');
            link.href = href;
            link.setAttribute('download', `resume.pdf`);
            document.body.appendChild(link);
            link.click();
            
            document.body.removeChild(link);
            URL.revokeObjectURL(href);
        }).catch(async error => {
            const isJsonBlob = (data) => data instanceof Blob && data.type === "application/json";
            const responseData = isJsonBlob(error?.response?.data) ? await (error?.response?.data)?.text() : error?.response?.data || {};
            const responseJson = (typeof responseData === "string") ? JSON.parse(responseData) : responseData;

            errorMessage.value = responseJson?.message || 'An error occurred while generating the resume.';
            console.error(error);
        }).finally(() => {
            isLoading.value = false;
        });
    }
</script>

<template>
    <div>
        <div v-if="errorMessage" class="p-4 mb-4 text-sm text-red-800 rounded-lg bg-red-50 dark:bg-gray-800 dark:text-red-400" role="alert">
            <span class="font-medium">Error!</span> {{ errorMessage }}
        </div>

        <div>
            <label for="name" class="block text-sm font-medium text-gray-700 mb-2">
                Resume.io Rendering Token
            </label>
            <input
                type="text"
                id="rendering_token"
                name="rendering_token"
                class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="Enter your rendering token"
                v-model="renderingToken"
                autocomplete="off"
            />
            <p class="my-2 text-xs text-gray-500">
                Please enter your resume.io rendering token.
            </p>
            <p class="my-2 text-xs text-gray-500">
                Don't know how to get it? Login to resume.io first then open this url <a class="text-blue-500" href="https://resume.io/api/app/resumes" target="_blank">https://resume.io/api/app/resumes</a>
            </p>
        </div>

        <div>
            <button
                type="submit"
                class="w-full mt-5 px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                @click="generate"
                :disabled="!renderingToken || renderingToken == '' || isLoading"
                >
                Generate
            </button>
        </div>
    </div>
</template>
