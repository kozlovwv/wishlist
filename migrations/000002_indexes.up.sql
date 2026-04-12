CREATE INDEX idx_wishlists_user_id ON wishlists(user_id);
CREATE INDEX idx_items_wishlist_id ON items(wishlist_id);
CREATE INDEX idx_wishlists_public_token ON wishlists(public_token);