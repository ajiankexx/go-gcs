local all_buf_ids = vim.api.nvim_list_bufs()
print("\nAll Buffer IDs:")
for _, buf_id in ipairs(all_buf_ids) do
    local buf_name = vim.api.nvim_buf_get_name(buf_id)
    local is_loaded = vim.api.nvim_buf_is_loaded(buf_id)
    local is_valid = vim.api.nvim_buf_is_valid(buf_id)
    print(string.format("  ID: %d, Name: '%s', Loaded: %s, Valid: %s",
                        buf_id, buf_name, tostring(is_loaded), tostring(is_valid)))
end
