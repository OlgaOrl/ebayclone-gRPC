#!/usr/bin/env python3
"""
eBayClone gRPC Python Client Example

This example demonstrates how to use the eBayClone gRPC API from Python.
It shows that the gRPC service is language-agnostic.

Prerequisites:
    pip install grpcio grpcio-tools

To generate Python protobuf code:
    python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. proto/ebayclone.proto
"""

import grpc
import sys
import os

# Add the proto directory to the path
sys.path.append(os.path.join(os.path.dirname(__file__), '..'))

try:
    import proto.ebayclone_pb2 as pb2
    import proto.ebayclone_pb2_grpc as pb2_grpc
except ImportError:
    print("Error: Proto files not generated for Python.")
    print("Run: python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. proto/ebayclone.proto")
    sys.exit(1)

def main():
    # Connect to gRPC server
    channel = grpc.insecure_channel('localhost:50051')
    
    # Create service stubs
    user_stub = pb2_grpc.UserServiceStub(channel)
    session_stub = pb2_grpc.SessionServiceStub(channel)
    listing_stub = pb2_grpc.ListingServiceStub(channel)
    order_stub = pb2_grpc.OrderServiceStub(channel)

    print("=== eBayClone gRPC Python Client Example ===")

    try:
        # 1. Create a user
        print("\n1. Creating user...")
        user = user_stub.CreateUser(pb2.UserCreate(
            username="pythonuser",
            email="python@example.com",
            password="pythonpass"
        ))
        print(f"Created user: ID={user.id}, Username={user.username}, Email={user.email}")

        # 2. Login
        print("\n2. Logging in...")
        login_resp = session_stub.Login(pb2.UserLogin(
            email="python@example.com",
            password="pythonpass"
        ))
        print(f"Login successful, token: {login_resp.token[:20]}...")

        # 3. Create a listing
        print("\n3. Creating listing...")
        listing = listing_stub.CreateListing(pb2.ListingCreate(
            title="MacBook Pro",
            description="Excellent condition laptop",
            price=1299.99,
            category="electronics",
            condition="like-new",
            location="San Francisco, CA"
        ))
        print(f"Created listing: ID={listing.id}, Title={listing.title}, Price=${listing.price}")

        # 4. Search listings
        print("\n4. Searching listings...")
        listings_resp = listing_stub.GetListings(pb2.ListingsRequest(
            search="MacBook",
            price_min=1000,
            price_max=2000
        ))
        print(f"Found {len(listings_resp.listings)} listings")

        # 5. Create an order
        print("\n5. Creating order...")
        order = order_stub.CreateOrder(pb2.OrderCreate(
            listing_id=listing.id,
            quantity=1,
            shipping_address=pb2.Address(
                street="456 Python St",
                city="San Francisco",
                state="CA",
                zip_code="94102",
                country="USA"
            ),
            buyer_notes="Handle with care"
        ))
        print(f"Created order: ID={order.id}, Status={order.status}, Total=${order.total_price}")

        # 6. Update order status
        print("\n6. Updating order status...")
        updated_order = order_stub.UpdateOrderStatus(pb2.UpdateOrderStatusRequest(
            id=order.id,
            status="shipped"
        ))
        print(f"Updated order status to: {updated_order.status}")

        print("\n=== Python client example completed successfully! ===")

    except grpc.RpcError as e:
        print(f"gRPC Error: {e.code()} - {e.details()}")
    except Exception as e:
        print(f"Error: {e}")
    finally:
        channel.close()

if __name__ == "__main__":
    main()
