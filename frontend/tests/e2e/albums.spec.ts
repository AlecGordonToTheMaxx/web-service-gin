import { expect, test } from '@playwright/test';

test.describe('Album Manager', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should display the Album Manager title', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Album Manager');
  });

  test('should create a new album', async ({ page }) => {
    // Fill out the form
    await page.fill('input[name="title"]', 'The Wall');
    await page.fill('input[name="artist"]', 'Pink Floyd');
    await page.fill('input[name="price"]', '24.99');

    // Submit the form
    await page.click('button[type="submit"]');

    // Wait for the album to appear in the list
    await expect(page.locator('text=The Wall')).toBeVisible({ timeout: 5000 });
    await expect(page.locator('text=Pink Floyd')).toBeVisible();
    await expect(page.locator('text=$24.99')).toBeVisible();
  });

  test('should edit an album', async ({ page }) => {
    // First create an album
    await page.fill('input[name="title"]', 'Dark Side of the Moon');
    await page.fill('input[name="artist"]', 'Pink Floyd');
    await page.fill('input[name="price"]', '19.99');
    await page.click('button[type="submit"]');

    // Wait for album to appear
    await expect(page.locator('text=Dark Side of the Moon')).toBeVisible({ timeout: 5000 });

    // Click edit button
    await page.locator('button:has-text("Edit")').first().click();

    // Form should be populated
    await expect(page.locator('input[name="title"]')).toHaveValue('Dark Side of the Moon');

    // Change the price
    await page.fill('input[name="price"]', '22.99');
    await page.click('button[type="submit"]');

    // Verify the updated price
    await expect(page.locator('text=$22.99')).toBeVisible({ timeout: 5000 });
  });

  test('should delete an album', async ({ page }) => {
    // Create an album
    await page.fill('input[name="title"]', 'Wish You Were Here');
    await page.fill('input[name="artist"]', 'Pink Floyd');
    await page.fill('input[name="price"]', '18.99');
    await page.click('button[type="submit"]');

    // Wait for album to appear
    await expect(page.locator('text=Wish You Were Here')).toBeVisible({ timeout: 5000 });

    // Setup dialog handler before clicking delete
    page.on('dialog', (dialog) => dialog.accept());

    // Click delete button
    await page.locator('button:has-text("Delete")').first().click();

    // Verify album is removed
    await expect(page.locator('text=Wish You Were Here')).not.toBeVisible({ timeout: 5000 });
  });

  test('should cancel editing', async ({ page }) => {
    // Create an album
    await page.fill('input[name="title"]', 'Animals');
    await page.fill('input[name="artist"]', 'Pink Floyd');
    await page.fill('input[name="price"]', '17.99');
    await page.click('button[type="submit"]');

    await expect(page.locator('text=Animals')).toBeVisible({ timeout: 5000 });

    // Click edit
    await page.locator('button:has-text("Edit")').first().click();

    // Form should show "Edit Album"
    await expect(page.locator('h2:has-text("Edit Album")')).toBeVisible();

    // Click cancel
    await page.click('button:has-text("Cancel")');

    // Form should reset to "Add New Album"
    await expect(page.locator('h2:has-text("Add New Album")')).toBeVisible();

    // Form fields should be empty
    await expect(page.locator('input[name="title"]')).toHaveValue('');
  });

  test('should show validation for required fields', async ({ page }) => {
    // Try to submit empty form
    await page.click('button[type="submit"]');

    // HTML5 validation should prevent submission
    const titleInput = page.locator('input[name="title"]');
    const isInvalid = await titleInput.evaluate((el: HTMLInputElement) => !el.validity.valid);
    expect(isInvalid).toBe(true);
  });

  test('should display album count', async ({ page }) => {
    // Check initial count (may vary depending on test isolation)
    const countText = await page.locator('p:has-text("Albums")').textContent();
    expect(countText).toMatch(/Albums \(\d+\)/);

    // Add an album
    await page.fill('input[name="title"]', 'The Division Bell');
    await page.fill('input[name="artist"]', 'Pink Floyd');
    await page.fill('input[name="price"]', '21.99');
    await page.click('button[type="submit"]');

    await expect(page.locator('text=The Division Bell')).toBeVisible({ timeout: 5000 });

    // Count should update
    const newCountText = await page.locator('p:has-text("Albums")').textContent();
    expect(newCountText).toMatch(/Albums \(\d+\)/);
  });
});
