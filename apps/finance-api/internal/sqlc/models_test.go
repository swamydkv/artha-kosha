package sqlc

import (
	"testing"
)


func TestAuditResult_ScanValue(t *testing.T) {
    var e AuditResult
    var ne NullAuditResult
    
    // Test Scan on base type
    if err := e.Scan(string("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := e.Scan([]byte("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := e.Scan(123); err == nil {
        t.Errorf("expected error")
    }
    
    // Test Scan on Null type
    if err := ne.Scan(nil); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := ne.Scan(string("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    
    // Test Value on Null type
    ne.Valid = false
    v1, _ := ne.Value()
    if v1 != nil {
        t.Errorf("expected nil")
    }
    
    ne.Valid = true
    ne.AuditResult = "valid"
    v2, _ := ne.Value()
    if v2 != string("valid") {
        t.Errorf("expected valid, got %v", v2)
    }
}

func TestDomainEventType_ScanValue(t *testing.T) {
    var e DomainEventType
    var ne NullDomainEventType
    
    // Test Scan on base type
    if err := e.Scan(string("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := e.Scan([]byte("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := e.Scan(123); err == nil {
        t.Errorf("expected error")
    }
    
    // Test Scan on Null type
    if err := ne.Scan(nil); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := ne.Scan(string("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    
    // Test Value on Null type
    ne.Valid = false
    v1, _ := ne.Value()
    if v1 != nil {
        t.Errorf("expected nil")
    }
    
    ne.Valid = true
    ne.DomainEventType = "valid"
    v2, _ := ne.Value()
    if v2 != string("valid") {
        t.Errorf("expected valid, got %v", v2)
    }
}

func TestProcessingStatus_ScanValue(t *testing.T) {
    var e ProcessingStatus
    var ne NullProcessingStatus
    
    // Test Scan on base type
    if err := e.Scan(string("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := e.Scan([]byte("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := e.Scan(123); err == nil {
        t.Errorf("expected error")
    }
    
    // Test Scan on Null type
    if err := ne.Scan(nil); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := ne.Scan(string("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    
    // Test Value on Null type
    ne.Valid = false
    v1, _ := ne.Value()
    if v1 != nil {
        t.Errorf("expected nil")
    }
    
    ne.Valid = true
    ne.ProcessingStatus = "valid"
    v2, _ := ne.Value()
    if v2 != string("valid") {
        t.Errorf("expected valid, got %v", v2)
    }
}

func TestSessionStatus_ScanValue(t *testing.T) {
    var e SessionStatus
    var ne NullSessionStatus
    
    // Test Scan on base type
    if err := e.Scan(string("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := e.Scan([]byte("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := e.Scan(123); err == nil {
        t.Errorf("expected error")
    }
    
    // Test Scan on Null type
    if err := ne.Scan(nil); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if err := ne.Scan(string("valid")); err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    
    // Test Value on Null type
    ne.Valid = false
    v1, _ := ne.Value()
    if v1 != nil {
        t.Errorf("expected nil")
    }
    
    ne.Valid = true
    ne.SessionStatus = "valid"
    v2, _ := ne.Value()
    if v2 != string("valid") {
        t.Errorf("expected valid, got %v", v2)
    }
}
